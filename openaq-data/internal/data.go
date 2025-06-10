package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"resty.dev/v3"
)

type openAQLocation struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Locality string `json:"locality"`
	Timezone string `json:"timezone"`
	Country  struct {
		Id   int32  `json:"id"`
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
}

type openAQLocationResponse struct {
	Results []openAQLocation `json:"results"`
}

type openAQMeasurement struct {
	Date struct {
		Utc   string `json:"utc"`
		Local string `json:"local"`
	} `json:"datetime"`
	Value       float64 `json:"value"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	Parameter struct {
		Id int32 `json:"id"`
	} `json:"parameter"`
	LocationId int32 `json:"locationId"`
}

type openAQMeasurementResponse struct {
	Results []openAQMeasurement `json:"results"`
}

type DataService struct {
	APIKey string
	client *resty.Client
	store  *Store
}

func NewDataService(apiKey, mongoURI string) (*DataService, error) {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")
	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}

	s, err := NewStore(mongoURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &DataService{
		APIKey: apiKey,
		client: client,
		store:  s,
	}, nil
}

func (s *DataService) Run(ctx context.Context) error {
	if err := s.FetchLocations(ctx); err != nil {
		return fmt.Errorf("initial locations fetch failed: %w", err)
	}

	if err := s.FetchMeasurements(ctx); err != nil {
		log.Printf("Failed to fetch measurements: %v", err)
	} else {
		log.Printf("Measurements updated at %s", time.Now().Format(time.RFC3339))
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.FetchMeasurements(ctx); err != nil {
				log.Printf("Failed to fetch measurements: %v", err)
			} else {
				log.Printf("Measurements updated at %s", time.Now().Format(time.RFC3339))
			}
		}
	}
}

func (s *DataService) FetchLocations(ctx context.Context) error {
	var data openAQLocationResponse
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"iso":   "PL",
			"limit": "1000",
		}).
		SetResult(&data).
		Get("locations")

	if err != nil {
		return fmt.Errorf("failed to fetch locations: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	for _, apiLoc := range data.Results {
		location := Location{
			ID:       apiLoc.Id,
			Name:     apiLoc.Name,
			Locality: apiLoc.Locality,
			Timezone: apiLoc.Timezone,
			Country: Country{
				ID:   apiLoc.Country.Id,
				Code: apiLoc.Country.Code,
				Name: apiLoc.Country.Name,
			},
			Sensors: []Sensor{},
		}
		for _, apiSensor := range apiLoc.Sensors {
			sensor := Sensor{
				ID:   apiSensor.Id,
				Name: apiSensor.Name,
				Parameter: Parameter{
					ID:          apiSensor.Parameter.Id,
					Name:        apiSensor.Parameter.Name,
					Units:       apiSensor.Parameter.Units,
					DisplayName: apiSensor.Parameter.DisplayName,
				},
			}
			location.Sensors = append(location.Sensors, sensor)
		}
		if err := s.store.StoreLocation(ctx, location); err != nil {
			log.Printf("Failed to upsert location %d: %v", location.ID, err)
		}
	}
	return nil
}

func (s *DataService) FetchMeasurements(ctx context.Context) error {
	locations, err := s.store.GetLocations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch locations from store: %w", err)
	}

	for _, loc := range locations {
		for _, sensor := range loc.Sensors {
			var apiData openAQMeasurementResponse
			resp, err := s.client.R().
				SetQueryParams(map[string]string{
					"parameter_id": fmt.Sprintf("%d", sensor.Parameter.ID),
				}).
				SetResult(&apiData).
				Get(fmt.Sprintf("locations/%d/latest", loc.ID))

			if err != nil {
				log.Printf("Failed to fetch measurements for location %s, sensor %s: %v", loc.Name, sensor.Name, err)
				continue
			}
			if resp.IsError() {
				log.Printf("API error for location %s, sensor %s: %s", loc.Name, sensor.Name, resp.Status())
				continue
			}

			for _, m := range apiData.Results {
				parsedTime, err := time.Parse(time.RFC3339Nano, m.Date.Utc)
				if err != nil {
					parsedTime, err = time.Parse(time.RFC3339, m.Date.Utc)
					if err != nil {
						log.Printf("Failed to parse time for measurement (UTC: %s): %v", m.Date.Utc, err)
						continue
					}
				}
				measurement := Measurement{
					DateTime: MeasurementDateTime{
						UTC:   m.Date.Utc,
						Local: m.Date.Local,
					},
					Timestamp: parsedTime,
					Value:     m.Value,
					Coordinates: Coordinates{
						Latitude:  m.Coordinates.Latitude,
						Longitude: m.Coordinates.Longitude,
					},
					SensorID:   sensor.ID,
					LocationID: loc.ID,
				}
				if err := s.store.StoreMeasurement(ctx, measurement); err != nil {
					log.Printf("Failed to insert measurement for location %s, sensor %s: %v", loc.Name, sensor.Name, err)
				}
			}
			<-time.After(1 * time.Second)
		}
	}
	return nil
}

func (s *DataService) Close() error {
	if s.store != nil {
		return s.store.Close()
	}
	return nil
}
