package internal

import (
	"fmt"

	"resty.dev/v3"
)

const (
	WrzeszczId    = 3346528
	NoweSzkotyId  = 3285962
	NowyPortId    = 3285963
	SrodmiescieId = 3359900
)

type reponse struct {
	Results []struct {
		Locality    string `json:"locality"`
		Name        string `json:"name"`
		Coordinates struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		Country struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"country"`
		Sensors []struct {
			Id        int32  `json:"id"`
			Name      string `json:"name"`
			Parameter struct {
				Id          int32  `json:"id"`
				Name        string `json:"name"`
				Units       string `json:"units"`
				DisplayName string `json:"displayName"`
			} `json:"parameter"`
		} `json:"sensors"`
	} `json:"results"`
}

type Service struct {
	APIKey string
	client *resty.Client
}

func NewService(apiKey string) *Service {
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

func (s *Service) FetchData(l string) (*reponse, error) {
	location, err := parseLocation(l)
	if err != nil {
		return nil, fmt.Errorf("invalid location: %w", err)
	}
	resp, err := s.client.R().
		SetResult(&reponse{}).
		Get("locations/" + fmt.Sprint(location))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	data := resp.Result().(*reponse)
	return data, nil
}

func parseLocation(location string) (int32, error) {
	switch location {
	case "wrzeszcz":
		return WrzeszczId, nil
	case "nowe-szkoty":
		return NoweSzkotyId, nil
	case "nowy-port":
		return NowyPortId, nil
	case "srodmiescie":
		return SrodmiescieId, nil
	default:
		return 0, fmt.Errorf(
			"unknown location: %s, possible values are: wrzeszcz, nowe-szkoty, nowy-port, srodmiescie",
			location,
		)
	}
}
