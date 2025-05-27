package internal

import (
	"fmt"

	"resty.dev/v3"
)

type locationResponse struct {
	Results []location `json:"results"`
}

type location struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Locality string `json:"locality"`
	Timezone string `json:"timezone"`
	Country  struct {
		Id   int32  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`
}

type measurementResponse struct {
	Results []measurement `json:"results"`
}
type measurement struct {
	DateTime struct {
		Utc   string `json:"utc"`
		Local string `json:"local"`
	} `json:"datetime"`
	Value       float64 `json:"value"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	SensorsId   int32 `json:"sensorsId"`
	LocationsId int32 `json:"locationsId"`
}

type DataService struct {
	APIKey       string
	client       *resty.Client
	locations    []location                 // TODO: store locations in database
	measurements map[location][]measurement // TODO: store measurements in database
}

func NewDataService(apiKey string) *DataService {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")

	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}
	return &DataService{
		APIKey:       apiKey,
		client:       client,
		locations:    []location{},
		measurements: make(map[location][]measurement),
	}
}

func (s *DataService) FetchLocations() error {
	var data locationResponse
	resp, err := s.client.R().
		SetQueryParam("iso", "PL").
		SetQueryParam("limit", "10").
		SetResult(&data).
		Get("locations")

	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	s.locations = data.Results
	return nil
}

func (s *DataService) FetchMeasurements() error {
	for _, loc := range s.locations {
		var data measurementResponse
		resp, err := s.client.R().
			SetResult(&data).
			Get(fmt.Sprintf("locations/%d/latest", loc.Id))

		if err != nil {
			return fmt.Errorf("failed to fetch measurements for location %s: %w", loc.Name, err)
		}
		if resp.IsError() {
			return fmt.Errorf("API error for location %s: %s", loc.Name, resp.Status())
		}

		s.measurements[loc] = data.Results
	}
	return nil
}
