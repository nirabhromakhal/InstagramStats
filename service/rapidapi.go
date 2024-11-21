package service

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type RapidApiService struct {
	host string
	key  string
}

func NewRapidApiService(host string, key string) *RapidApiService {
	return &RapidApiService{host: host, key: key}
}

func (r *RapidApiService) GetData(BaseUrl string, params map[string]string) ([]byte, error) {
	query := url.Values{}
	for k, v := range params {
		query.Add(k, v)
	}
	CompleteUrl := BaseUrl + "?" + query.Encode()
	req, err := http.NewRequest("GET", CompleteUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-host", r.host)
	req.Header.Add("x-rapidapi-key", r.key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	}
	return nil, errors.New(resp.Status)
}
