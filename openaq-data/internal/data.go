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

type dataResponse struct {
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

type DataService struct {
	APIKey          string
	client          *resty.Client
	locationService *LocationService
}

func NewDataService(apiKey string, ls *LocationService) *DataService {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")

	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}
	return &DataService{
		APIKey:          apiKey,
		client:          client,
		locationService: ls,
	}
}

func (s *DataService) FetchData(l string) (*dataResponse, error) {
	coords, err := s.locationService.FetchLocation(l)
	if err != nil {
		return nil, fmt.Errorf("error fetching location: %w", err)
	}
	resp, err := s.client.R().
		SetQueryParam("coordinates", fmt.Sprintf("%s,%s", coords.Lat, coords.Lon)).
		SetQueryParam("radius", "10000").
		SetQueryParam("limit", "10").
		SetResult(&dataResponse{}).
		Get("locations")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	data := resp.Result().(*dataResponse)
	return data, nil
}

func (s *DataService) FetchLocationCoords(location string) (*locationResponse, error) {
	resp, err := s.locationService.FetchLocation(location)
	if err != nil {
		return nil, fmt.Errorf("error fetching location: %w", err)
	}
	return resp, nil
}
