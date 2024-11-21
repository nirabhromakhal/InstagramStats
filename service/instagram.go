package service

import (
	"InstagramStats/entity"
	"InstagramStats/enum"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type InstagramService struct {
	db              *gorm.DB
	rapidApiService *RapidApiService
}

type InstagramProfileResponse struct {
	Data struct {
		FullName        string `json:"full_name"`
		Username        string `json:"username"`
		FollowerCount   int64  `json:"follower_count"`
		FollowingCount  int64  `json:"following_count"`
		Id              string `json:"id"`
		ProfileImageUrl string `json:"profile_pic_url"`
		IsPrivate       bool   `json:"is_private"`
	} `json:"data"`
}

type InstagramVideosResponse struct {
	Data struct {
		Count int64 `json:"count"`
		Items []struct {
			Id           string `json:"id"`
			LikeCount    int64  `json:"like_count"`
			CommentCount int64  `json:"comment_count"`
			ThumbnailUrl string `json:"thumbnail_url"`
			VideoUrl     string `json:"video_url"`
			TakenAt      int64  `json:"taken_at"`
			Caption      struct {
				Text string `json:"text"`
			} `json:"caption"`
		} `json:"items"`
	} `json:"data"`
	PaginationToken string `json:"pagination_token"`
}

func NewInstagramService(db *gorm.DB, service *RapidApiService) *InstagramService {
	return &InstagramService{db: db, rapidApiService: service}
}

func (service *InstagramService) AddInstagramChannel(username string) error {
	resp, err := service.getInstagramProfile(username)
	if err != nil {
		return err
	}
	// fetch user from db
	var user entity.User
	service.db.First(&user, "LOWER(username) = LOWER(?)", username)

	// link user to channel
	service.ingestChannelIntoDatabase(user, resp)
	return nil
}

func (service *InstagramService) getInstagramProfile(username string) (*InstagramProfileResponse, error) {
	baseUrl := "https://instagram-scraper-api2.p.rapidapi.com/v1/info"
	queryParams := map[string]string{
		"username_or_id_or_url": username,
	}
	data, err := service.rapidApiService.GetData(baseUrl, queryParams)
	if err != nil {
		return nil, err
	}

	var response InstagramProfileResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &response, nil
}

func (service *InstagramService) ingestChannelIntoDatabase(user entity.User, accountInfo *InstagramProfileResponse) {
	channel := entity.Channel{
		Name:            accountInfo.Data.Username,
		ChannelType:     enum.ChannelTypeInstagram,
		VendorChannelId: accountInfo.Data.Id,
		Thumbnail:       accountInfo.Data.ProfileImageUrl,
		ChannelMembers: []entity.ChannelMember{
			{
				User:           user,
				UserID:         user.ID,
				Role:           enum.ChannelMemberRoleOwner,
				Visibility:     enum.ChannelMemberVisibilityDefault,
				ProfilePicture: accountInfo.Data.ProfileImageUrl,
			},
		},
	}
	service.db.Create(&channel)
	fmt.Printf("Successfully created channel %s in db\n", channel.Name)
}

func (service *InstagramService) updateChannelIntoDatabase(channel *entity.Channel, accountInfo *InstagramProfileResponse) {
	channel.Name = accountInfo.Data.Username
	channel.Thumbnail = accountInfo.Data.ProfileImageUrl
	channel.ChannelMembers[0].ProfilePicture = accountInfo.Data.ProfileImageUrl
	service.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(channel)
	fmt.Printf("Successfully updated channel %s in db\n", channel.Name)
}

func (service *InstagramService) AddInstagramPostsReels(username string) error {
	resp, err := service.getInstagramPostsReels(username)
	if err != nil {
		return err
	}

	// fetch channel from db
	var channel entity.Channel
	service.db.First(&channel, "LOWER(name) = LOWER(?)", username)

	// link videos to channel
	service.ingestVideosIntoDatabase(&channel, resp)
	return nil
}

func (service *InstagramService) getInstagramPostsReels(username string) (*InstagramVideosResponse, error) {
	baseUrl := "https://instagram-scraper-api2.p.rapidapi.com/v1/posts"
	queryParams := map[string]string{
		"username_or_id_or_url": username,
	}
	data, err := service.rapidApiService.GetData(baseUrl, queryParams)
	if err != nil {
		return nil, err
	}

	var aggregatedResponse InstagramVideosResponse
	if err := json.Unmarshal(data, &aggregatedResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// keep querying to get all posts
	for aggregatedResponse.PaginationToken != "" {
		queryParams["pagination_token"] = aggregatedResponse.PaginationToken
		data, err = service.rapidApiService.GetData(baseUrl, queryParams)
		if err != nil {
			return nil, err
		}
		var response InstagramVideosResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %v", err)
		}
		aggregatedResponse.Data.Count += response.Data.Count
		aggregatedResponse.Data.Items = append(aggregatedResponse.Data.Items, response.Data.Items...)
		aggregatedResponse.PaginationToken = response.PaginationToken
	}
	return &aggregatedResponse, nil
}

func (service *InstagramService) ingestVideosIntoDatabase(channel *entity.Channel, resp *InstagramVideosResponse) {
	var existingVideos []entity.Video
	service.db.Where("channel_id = ?", channel.ID).Find(&existingVideos)

	var videos []entity.Video
	for _, item := range resp.Data.Items {
		// check if this item already exists in db
		exists := false
		for _, existingVideo := range existingVideos {
			if existingVideo.VideoID == item.Id {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		// if this item does not exist, ingest it
		videos = append(videos, entity.Video{
			VideoID:     item.Id,
			ChannelID:   channel.ID,
			ChannelName: channel.Name,
			Title:       item.Caption.Text,
			PublishedAt: time.Unix(item.TakenAt, 0),
			Thumbnail:   item.ThumbnailUrl,
			VideoURL:    item.VideoUrl,
			Platform:    enum.ChannelTypeInstagram,
		})
	}
	if videos != nil {
		service.db.Create(&videos)
	}
	fmt.Printf("Successfully ingested videos for channel %s\n", channel.Name)
}

func (service *InstagramService) UpdateAllInstagramChannels() error {
	channels, err := service.getAllInstagramChannels()
	if err != nil {
		return err
	}
	for _, channel := range channels {
		// update profile info
		profileResp, err := service.getInstagramProfile(channel.Name)
		if err != nil {
			fmt.Printf("Failed to get instagram profile for %s: %v\n", channel.Name, err)
		} else {
			service.updateChannelIntoDatabase(&channel, profileResp)
		}

		// update videos for public account
		if profileResp.Data.IsPrivate {
			fmt.Printf("Skipping video ingestion for private Instagram account %s\n", channel.Name)
			continue
		}
		postResp, err := service.getInstagramPostsReels(channel.Name)
		if err != nil {
			fmt.Printf("Failed to get instagram posts for %s: %v\n", channel.Name, err)
		} else {
			service.ingestVideosIntoDatabase(&channel, postResp)
		}
	}
	return nil
}

func (service *InstagramService) getAllInstagramChannels() ([]entity.Channel, error) {
	var channels []entity.Channel
	err := service.db.Preload("ChannelMembers").Where("channel_type = ?", enum.ChannelTypeInstagram).Find(&channels).Error
	return channels, err
}
