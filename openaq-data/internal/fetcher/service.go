package fetcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"openaq-data/internal/fetcher/api"
	"openaq-data/internal/store"
	"openaq-data/types"
	"sync"
	"time"
)

const (
	stationsUpdateInterval     = 24 * time.Hour
	measurementsUpdateInterval = 1 * time.Hour
)

var (
	ErrStationsNotLoaded = errors.New("stations not loaded yet")
)

// Fetcher service is responsible for fetching data from the OpenAQ API
// and storing it in the local store.
// It runs in the background, periodically updating stations and measurements.
type Service struct {
	api    *api.Service
	store  *store.Store
	logger *slog.Logger

	stationsLoadedOnce sync.Once
	stationsLoaded     chan struct{}
}

func NewService(apiKey string, s *store.Store, l *slog.Logger) (*Service, error) {
	return &Service{
		api:            api.New(apiKey, l),
		store:          s,
		logger:         l,
		stationsLoaded: make(chan struct{}, 1),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	go s.updateStationsLoop(ctx)
	go s.updateMeasurementsLoop(ctx)

	<-ctx.Done()
	s.logger.Info("Service is shutting down...")
	return nil
}

func (s *Service) updateStationsLoop(ctx context.Context) {
	ticker := time.NewTicker(stationsUpdateInterval)
	defer ticker.Stop()

	for {
		if err := s.loadStations(ctx); err != nil {
			s.logger.Error("Failed to update stations", slog.Any("error", err))
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) loadStations(ctx context.Context) error {
	locations, err := s.api.FetchLocations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch locations from API: %w", err)
	}

	stations := s.buildStationsFromApiData(locations)
	if err := s.store.StoreStations(ctx, stations); err != nil {
		return fmt.Errorf("failed to store stations: %w", err)
	}

	s.stationsLoadedOnce.Do(func() {
		close(s.stationsLoaded)
	})

	return nil
}

func (s *Service) buildStationsFromApiData(locations []api.OpenAQLocation) []types.Station {
	var stations []types.Station
	for _, apiLoc := range locations {
		station := types.Station{
			Id:       &apiLoc.Id,
			Name:     &apiLoc.Name,
			Locality: &apiLoc.Locality,
			Timezone: &apiLoc.Timezone,
			Country: &types.Country{
				Id:   &apiLoc.Country.Id,
				Code: &apiLoc.Country.Code,
				Name: &apiLoc.Country.Name,
			},
			Coordinates: &types.Coordinates{
				Latitude:  &apiLoc.Coordinates.Latitude,
				Longitude: &apiLoc.Coordinates.Longitude,
			},
			Parameters: &[]types.Parameter{},
		}
		for _, apiSensor := range apiLoc.Sensors {
			parameter := &types.Parameter{
				Id:          &apiSensor.Parameter.Id,
				Name:        &apiSensor.Parameter.Name,
				Units:       &apiSensor.Parameter.Units,
				DisplayName: &apiSensor.Parameter.DisplayName,
			}
			*station.Parameters = append(*station.Parameters, *parameter)
		}
		stations = append(stations, station)
	}
	return stations
}

func (s *Service) updateMeasurementsLoop(ctx context.Context) {
	ticker := time.NewTicker(measurementsUpdateInterval)
	defer ticker.Stop()

	for {
		if err := s.updateMeasurements(ctx); err != nil {
			s.logger.Error("Failed to update measurements", slog.Any("error", err))
			if errors.Is(err, ErrStationsNotLoaded) {
				<-time.After(5 * time.Second)
				continue
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateMeasurements(ctx context.Context) error {
	select {
	case <-s.stationsLoaded:
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrStationsNotLoaded
	}

	stations, err := s.store.GetStations(ctx)
	if err != nil {
		return err
	}
	if len(stations) == 0 {
		return fmt.Errorf("no stations found in store")
	}
	for _, station := range stations {
		if err := s.updateMeasurementsForStation(ctx, station); err != nil {
			s.logger.Error(
				"Failed to update measurements for station",
				slog.Any("error", err),
			)
		}
	}
	return nil
}

func (s *Service) updateMeasurementsForStation(ctx context.Context, station types.Station) error {
	measurements, err := s.loadMeasurementsForStation(station)
	if err != nil {
		return err
	}
	if len(measurements) == 0 {
		return fmt.Errorf("no measurements found for station %s", *station.Name)
	}

	err = s.store.DeleteMeasurementsForStation(ctx, *station.Id)
	if err != nil {
		return fmt.Errorf("failed to delete existing measurements for station %s: %w", *station.Name, err)
	}

	err = s.store.StoreMeasurements(ctx, measurements)
	if err != nil {
		return fmt.Errorf("failed to store measurements for station %s: %w", *station.Name, err)
	}
	return nil
}

func (s *Service) loadMeasurementsForStation(station types.Station) ([]types.Measurement, error) {
	var measurements []types.Measurement
	for _, parameter := range *station.Parameters {
		apiData, err := s.tryGetMeasurementsForStation(*station.Id, *parameter.Id)
		if err != nil {
			return nil, err
		}

		paramMeasurements := s.buildMeasurementsFromApiData(apiData, &parameter, *station.Id)
		measurements = append(measurements, paramMeasurements...)

		// To avoid hitting rate limits
		<-time.After(1 * time.Second)
	}
	return measurements, nil
}

func (s *Service) tryGetMeasurementsForStation(stationId, paramId int32) ([]api.OpenAQMeasurement, error) {
	for range 3 {
		apiData, err := s.api.FetchMeasurementsForLocation(stationId, paramId)
		if err != nil {
			if errors.Is(err, api.ErrRateLimitExceeded) {
				s.logger.Warn("Rate limit exceeded, retrying after delay")
				<-time.After(5 * time.Second)
				continue
			}
			return nil, err
		}
		return apiData, nil
	}
	return nil, fmt.Errorf("failed to fetch measurements for station %d, parameter %d after retries", stationId, paramId)
}

func (s *Service) buildMeasurementsFromApiData(
	apiMeasurements []api.OpenAQMeasurement,
	param *types.Parameter,
	stationId int32,
) []types.Measurement {
	var measurements []types.Measurement
	for _, m := range apiMeasurements {
		parsedTime, err := time.Parse(time.RFC3339Nano, m.Date.Utc)
		if err != nil {
			parsedTime, err = time.Parse(time.RFC3339, m.Date.Utc)
			if err != nil {
				log.Printf("Failed to parse time for measurement (UTC: %s): %v", m.Date.Utc, err)
				continue
			}
		}
		measurements = append(measurements, types.Measurement{
			Datetime: &types.MeasurementDateTime{
				Utc:   &m.Date.Utc,
				Local: &m.Date.Local,
			},
			Timestamp: &parsedTime,
			Value:     &m.Value,
			Coordinates: &types.Coordinates{
				Latitude:  &m.Coordinates.Latitude,
				Longitude: &m.Coordinates.Longitude,
			},
			SensorId:  param.Id,
			StationId: &stationId,
		})
	}
	return measurements
}
