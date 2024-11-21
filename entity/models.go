package entity

import (
	"strings"
	"time"

	"InstagramStats/enum"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = strings.ToLower(ulid.Make().String())
	}

	return nil
}

type User struct {
	BaseModel
	Username           string `gorm:"uniqueIndex;not null"`
	Email              string `gorm:"uniqueIndex;not null"`
	PasswordHash       string
	Name               string
	DateOfBirth        string
	ProfilePicture     string
	About              string
	UserRole           enum.UserRole `gorm:"type:text;default:'CREATOR'"`
	UserNumber         int
	UserState          enum.UserState `gorm:"type:text;default:'WAITLISTED'"`
	SubscriptionPlanID string
	LastActiveAt       time.Time
}

func (u *User) SetDefaultPictureIfNotPresent() {
	if u.ProfilePicture == "" {
		u.ProfilePicture = "https://ui-avatars.com/api/?name=P&background=0096E5&color=FFFFFF"
	}
}

type Channel struct {
	BaseModel
	Name            string
	ChannelType     enum.ChannelType
	VendorChannelId string
	Thumbnail       string
	ChannelMembers  []ChannelMember
	Videos          []Video
}

func (c *Channel) HasMember(userID string, role enum.ChannelMemberRole) bool {
	for _, member := range c.ChannelMembers {
		if member.UserID == userID && member.Role == role {
			return true
		}
	}
	return false
}

func (ch *Channel) IsOwner(userID string) bool {
	for _, member := range ch.ChannelMembers {
		if member.UserID == userID && member.Role == enum.ChannelMemberRoleOwner {
			return true
		}
	}
	return false
}

type ChannelMember struct {
	BaseModel
	ChannelID      string
	Channel        Channel
	UserID         string
	User           User
	Role           enum.ChannelMemberRole
	Visibility     enum.ChannelMemberVisibility
	ProfilePicture string
}

type Video struct {
	BaseModel
	VideoID     string `gorm:"uniqueIndex;not null"`
	ChannelID   string `gorm:"index;not null"`
	ChannelName string
	Title       string
	PublishedAt time.Time
	Thumbnail   string
	VideoURL    string
	Platform    enum.ChannelType
}
