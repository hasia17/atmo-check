package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"resty.dev/v3"
)

type openAQLocation struct {
	Id       int32  `json:"id"       bson:"_id"`
	Name     string `json:"name"     bson:"name"`
	Locality string `json:"locality" bson:"locality"`
	Timezone string `json:"timezone" bson:"timezone"`
	Country  struct {
		Id   int32  `json:"id" bson:"id"`
		Code string `json:"code" bson:"code"`
		Name string `json:"name" bson:"name"`
	} `json:"country"  bson:"country"`
	Sensors []struct {
		Id        int32  `json:"id" bson:"_id"`
		Name      string `json:"name" bson:"name"`
		Parameter struct {
			Id          int32  `json:"id" bson:"id"`
			Name        string `json:"name" bson:"name"`
			Units       string `json:"units" bson:"units"`
			DisplayName string `json:"displayName" bson:"displayName"`
		} `json:"parameter" bson:"parameter"`
	}
}

type openAQLocationResponse struct {
	Results []openAQLocation `json:"results"`
}

type openAQMeasurement struct {
	DateTime struct {
		Utc   string `json:"utc" bson:"utc"`
		Local string `json:"local" bson:"local"`
	} `json:"datetime"    bson:"datetime"`
	Timestamp   time.Time `json:"-"           bson:"timestamp"`
	Value       float64   `json:"value"       bson:"value"`
	Coordinates struct {
		Latitude  float64 `json:"latitude" bson:"latitude"`
		Longitude float64 `json:"longitude" bson:"longitude"`
	} `json:"coordinates" bson:"coordinates"`
	SensorsId   int32 `json:"sensorsId"   bson:"sensorsId"`
	LocationsId int32 `json:"locationsId" bson:"locationsId"`
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

	for _, loc := range data.Results {
		if err := s.store.StoreLocation(ctx, loc); err != nil {
			log.Printf("Failed to upsert location %d: %v", loc.Id, err)
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
		var data openAQMeasurementResponse
		resp, err := s.client.R().
			SetResult(&data).
			Get(fmt.Sprintf("locations/%d/latest", loc.Id))
		if err != nil {
			log.Printf("Failed to fetch measurements for location %s: %v", loc.Name, err)
			continue
		}
		if resp.IsError() {
			log.Printf("API error for location %s: %s", loc.Name, resp.Status())
			continue
		}

		for _, m := range data.Results {
			if err := s.store.StoreMeasurement(ctx, m); err != nil {
				log.Printf("Failed to insert measurement for location %s: %v", loc.Name, err)
			}
		}
		// openaq rate limit is 60 requests per minute
		<-time.After(1 * time.Second)
	}
	return nil
}

func (s *DataService) Close() error {
	if s.store != nil {
		return s.store.Close()
	}
	return nil
}
