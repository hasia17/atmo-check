package internal

import (
	"fmt"

	"resty.dev/v3"
)

type locationResponse struct {
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Name string `json:"name"`
}

type LocationService struct {
	client *resty.Client
}

func NewLocationService() *LocationService {
	client := resty.New().
		SetBaseURL("https://nominatim.openstreetmap.org/search").
		SetHeader("Content-Type", "application/json").
		SetQueryParam("format", "json")
	return &LocationService{
		client: client,
	}
}

func (s *LocationService) FetchLocation(location string) (*locationResponse, error) {
	resp, err := s.client.R().
		SetQueryParam("q", location).
		SetResult(&locationResponse{}).
		Get("")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error fetching location: %s", resp.Status())
	}

	locationData := resp.Result().(*locationResponse)
	return locationData, nil
}
