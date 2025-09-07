package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"openaq-data/api"
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
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
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
	logger *slog.Logger
}

func NewDataService(apiKey, mongoURI string, l *slog.Logger) (*DataService, error) {
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
		logger: l,
	}, nil
}

func (s *DataService) Run(ctx context.Context) error {
	s.logger.Info("Starting DataService...")

	hasData, err := s.store.HasStations(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for existing data: %w", err)
	}
	if !hasData {
		s.logger.Info("No data found in DB. Running initial fetch...")
		if err := s.FetchLocations(ctx); err != nil {
			s.logger.Error("Initial locations fetch failed", slog.Any("error", err))
		}
		if err := s.FetchMeasurements(ctx); err != nil {
			s.logger.Error("Initial measurements fetch failed", slog.Any("error", err))
		}
	} else {
		s.logger.Info("Data found in DB. Skipping initial fetch.")
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.logger.Info("Fetching locations and measurements...")
			if err := s.FetchMeasurements(ctx); err != nil {
				log.Printf("Failed to fetch measurements: %v", err)
			} else {
				log.Printf("Measurements updated at %s", time.Now().Format(time.RFC3339))
			}
		}
	}
}

func (s *DataService) FetchLocations(ctx context.Context) error {
	s.logger.Info("Fetching locations from OpenAQ API...")
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
		station := api.Station{
			Id:       &apiLoc.Id,
			Name:     &apiLoc.Name,
			Locality: &apiLoc.Locality,
			Timezone: &apiLoc.Timezone,
			Country: &api.Country{
				Id:   &apiLoc.Country.Id,
				Code: &apiLoc.Country.Code,
				Name: &apiLoc.Country.Name,
			},
			Coordinates: &api.Coordinates{
				Latitude:  &apiLoc.Coordinates.Latitude,
				Longitude: &apiLoc.Coordinates.Longitude,
			},
			Parameters: &[]api.Parameter{},
		}
		for _, apiSensor := range apiLoc.Sensors {
			parameter := &api.Parameter{
				Id:          &apiSensor.Parameter.Id,
				Name:        &apiSensor.Parameter.Name,
				Units:       &apiSensor.Parameter.Units,
				DisplayName: &apiSensor.Parameter.DisplayName,
			}
			*station.Parameters = append(*station.Parameters, *parameter)
		}
		if err := s.store.StoreStation(ctx, station); err != nil {
			s.logger.Error("Failed to store location", slog.String("location", *station.Name), slog.Any("error", err))
		}
	}
	s.logger.Info("Stations fetched and stored successfully", slog.Int("count", len(data.Results)))
	return nil
}

func (s *DataService) FetchMeasurements(ctx context.Context) error {
	s.logger.Info("Fetching measurements from OpenAQ API...")
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch stations from store: %w", err)
	}
	if len(stations) == 0 {
		s.logger.Info("No stations found in store, skipping measurements fetch")
		return nil
	}

	for _, st := range stations {
		for _, parameter := range *st.Parameters {
			var apiData openAQMeasurementResponse
			resp, err := s.client.R().
				SetQueryParams(map[string]string{
					"parameter_id": fmt.Sprintf("%d", parameter.Id),
				}).
				SetResult(&apiData).
				Get(fmt.Sprintf("locations/%d/latest", *st.Id))

			if err != nil {
				log.Printf("Failed to fetch measurements for station %s, parameter %s: %v", st.Name, parameter.Name, err)
				continue
			}
			if resp.IsError() {
				log.Printf("API error for station %s, parameter %s: %s", st.Name, parameter.Name, resp.Status())
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
				measurement := api.Measurement{
					Datetime: &api.MeasurementDateTime{
						Utc:   &m.Date.Utc,
						Local: &m.Date.Local,
					},
					Timestamp: &parsedTime,
					Value:     &m.Value,
					Coordinates: &api.Coordinates{
						Latitude:  &m.Coordinates.Latitude,
						Longitude: &m.Coordinates.Longitude,
					},
					SensorId:  parameter.Id,
					StationId: st.Id,
				}
				if err := s.store.StoreMeasurement(ctx, measurement); err != nil {
					s.logger.Error("Failed to store measurement",
						slog.String("station", *st.Name),
						slog.String("parameter", *parameter.Name),
						slog.Any("measurement", measurement),
						slog.Any("error", err))
				}
			}
			<-time.After(1 * time.Second)
		}
	}
	return nil
}

func (s *DataService) SaveStation(ctx context.Context, station api.Station) error {
	s.logger.Info("Saving station", slog.String("station", *station.Name))
	err := s.store.StoreStation(ctx, station)
	if err != nil {
		s.logger.Error("Failed to save station", slog.String("station", *station.Name), slog.Any("error", err))
	}
	return err
}

func (s *DataService) SaveMeasurements(ctx context.Context, measurements []api.Measurement) error {
	for _, m := range measurements {
		s.logger.Info("Saving measurement", slog.Any("measurement", m))
		if err := s.store.StoreMeasurement(ctx, m); err != nil {
			s.logger.Error("Failed to save measurement", slog.Any("measurement", m), slog.Any("error", err))
			return err
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

func (s *DataService) GetStations(ctx context.Context) ([]api.Station, error) {
	s.logger.Info("Getting all stations")
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		s.logger.Error("Failed to get stations", slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("Stations fetched", slog.Int("count", len(stations)))
	return stations, nil
}

func (s *DataService) GetStationByID(ctx context.Context, id int32) (*api.Station, error) {
	s.logger.Info("Getting station by ID", slog.Int("stationId", int(id)))
	station, err := s.store.GetStationByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get station by ID", slog.Int("stationId", int(id)), slog.Any("error", err))
		return nil, err
	}
	if station == nil {
		s.logger.Info("Station not found", slog.Int("stationId", int(id)))
	} else {
		s.logger.Info("Station fetched", slog.Int("stationId", int(id)))
	}
	return station, nil
}

func (s *DataService) GetMeasurementsByStation(
	ctx context.Context,
	stationID int32,
	limit int64,
) ([]api.Measurement, error) {
	s.logger.Info("Getting measurements for station", slog.Int("stationId", int(stationID)), slog.Int64("limit", limit))
	measurements, err := s.store.GetMeasurementsByStation(ctx, stationID, limit)
	if err != nil {
		s.logger.Error(
			"Failed to get measurements for station",
			slog.Int("stationId", int(stationID)),
			slog.Any("error", err),
		)
		return nil, err
	}
	s.logger.Info("Measurements fetched", slog.Int("stationId", int(stationID)), slog.Int("count", len(measurements)))
	return measurements, nil
}

func (s *DataService) GetParametersByStationID(ctx context.Context, id int32) ([]api.Parameter, error) {
	s.logger.Info("Getting parameters for station", slog.Int("stationId", int(id)))
	parameters, err := s.store.GetParametersByStationID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get parameters for station", slog.Int("stationId", int(id)), slog.Any("error", err))
		return nil, err
	}
	if parameters == nil {
		s.logger.Info("No parameters found for station", slog.Int("stationId", int(id)))
	} else {
		s.logger.Info("Parameters fetched for station", slog.Int("stationId", int(id)), slog.Int("count", len(parameters)))
	}
	return parameters, nil
}

func (s *DataService) GetLatestMeasurementsByStation(ctx context.Context, stationID int32) ([]api.Measurement, error) {
	s.logger.Info("Getting latest measurements for station", slog.Int("stationId", int(stationID)))
	measurements, err := s.store.GetLatestMeasurementsByStation(ctx, stationID)
	if err != nil {
		s.logger.Error(
			"Failed to get latest measurements for station",
			slog.Int("stationId", int(stationID)),
			slog.Any("error", err),
		)
		return nil, err
	}
	s.logger.Info("Latest measurements fetched", slog.Int("stationId", int(stationID)), slog.Int("count", len(measurements)))
	return measurements, nil
}
