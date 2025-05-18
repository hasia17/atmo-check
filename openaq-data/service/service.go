package service

import (
	"fmt"

	"resty.dev/v3"
)

type reponse struct {
	Results []struct {
		Locality    string `json:"locality"`
		Name        string `json:"name"`
		Coordinates struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		Sensors []struct {
			Id   int32  `json:"id"`
			Name string `json:"name"`
		} `json:"sensors"`
	} `json:"results"`
}

type Service struct {
	APIKey string
	client *resty.Client
}

func New(apiKey string) *Service {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")

	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}
	return &Service{
		APIKey: apiKey,
		client: client,
	}
}

func (s *Service) FetchData(city string) (*reponse, error) {
	resp, err := s.client.R().
		SetResult(&reponse{}).
		Get("locations/2178") // TODO: Replace with dynamic endpoint based on city

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	data := resp.Result().(*reponse)
	return data, nil
}
