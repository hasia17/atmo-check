package data

import (
	"context"
	"fmt"
	"log/slog"
	"openaq-data/api"
	"openaq-data/internal/store"

	"resty.dev/v3"
)

type Service struct {
	APIKey string
	client *resty.Client
	store  *store.Store
	logger *slog.Logger
}

func NewService(apiKey, mongoURI string, l *slog.Logger) (*Service, error) {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")
	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}

	s, err := store.New(mongoURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &Service{
		APIKey: apiKey,
		client: client,
		store:  s,
		logger: l,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("Starting Service...")

	hasData, err := s.store.HasStations(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for existing data: %w", err)
	}
	if !hasData {
		s.logger.Info("No data found in DB. Running initial fetch...")
		if err := s.fetchLocations(ctx); err != nil {
			s.logger.Error("Initial locations fetch failed", slog.Any("error", err))
		}
	} else {
		s.logger.Info("Data found in DB. Skipping initial fetch.")
	}

	go s.updateMeasurementsLoop(ctx)

	<-ctx.Done()
	s.logger.Info("Service is shutting down...")
	return nil
}

func (s *Service) fetchLocations(ctx context.Context) error {
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

func (s *Service) Close() error {
	if s.store != nil {
		return s.store.Close()
	}
	return nil
}

func (s *Service) GetStations(ctx context.Context) ([]api.Station, error) {
	s.logger.Info("Getting all stations")
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		s.logger.Error("Failed to get stations", slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("Stations returned", slog.Int("count", len(stations)))
	return stations, nil
}

func (s *Service) GetStationByID(ctx context.Context, id int32) (*api.Station, error) {
	s.logger.Info("Getting station by ID", slog.Int("stationId", int(id)))
	station, err := s.store.GetStationByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get station by ID", slog.Int("stationId", int(id)), slog.Any("error", err))
		return nil, err
	}
	if station == nil {
		s.logger.Info("Station not found", slog.Int("stationId", int(id)))
	} else {
		s.logger.Info("Station returned", slog.Int("stationId", int(id)))
	}
	return station, nil
}

func (s *Service) GetMeasurementsForStation(
	ctx context.Context,
	stationID int32,
	limit int64,
) ([]api.Measurement, error) {
	station, err := s.store.GetStationByID(ctx, stationID)
	if err != nil {
		s.logger.Error("Failed to get station by ID", slog.Int("stationId", int(stationID)), slog.Any("error", err))
		return nil, err
	}
	if station == nil {
		s.logger.Info("Station not found", slog.Int("stationId", int(stationID)))
		return nil, fmt.Errorf("station with ID %d not found", stationID)
	}
	measurements, err := s.store.GetLatestMeasurementsByStation(ctx, *station.Id)
	if err != nil {
		s.logger.Error(
			"Failed to get measurements for station",
			slog.Int("stationId", int(stationID)),
			slog.Any("error", err),
		)
		return nil, err
	}

	if len(measurements) == 0 {
		newMeasurements, err := s.fetchMeasurementsForStation(*station)
		if err != nil {
			s.logger.Error(
				"Failed to fetch measurements for station",
				slog.Int("stationId", int(stationID)),
				slog.Any("error", err),
			)
			return nil, err
		}
		if err := s.store.DeleteMeasurementsForStation(ctx, *station.Id); err != nil {
			s.logger.Error(
				"Failed to delete existing measurements for station",
				slog.String("station", *station.Name),
				slog.Any("error", err),
			)
		}
		if err := s.store.StoreMeasurements(ctx, newMeasurements); err != nil {
			s.logger.Error(
				"Failed to store measurements for station",
				slog.String("station", *station.Name),
				slog.Any("error", err),
			)
		}
		measurements = newMeasurements
	}

	s.logger.Info("Measurements returned", slog.Int("stationId", int(stationID)), slog.Int("count", len(measurements)))
	return measurements, nil
}
