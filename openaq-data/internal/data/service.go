package data

import (
	"context"
	"fmt"
	"log/slog"
	"openaq-data/internal"
	"openaq-data/internal/store"
	"openaq-data/internal/types"
)

type Service struct {
	store  *store.Store
	logger *slog.Logger
}

func NewService(s *store.Store, l *slog.Logger) internal.DataService {
	return &Service{
		store:  s,
		logger: l,
	}
}

func (s *Service) Stations(ctx context.Context) ([]types.Station, error) {
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func (s *Service) Parameters(ctx context.Context) ([]types.Parameter, error) {
	parameters, err := s.store.GetParameters(ctx)
	if err != nil {
		return nil, err
	}
	return parameters, nil
}

func (s *Service) MeasurementsForStation(
	ctx context.Context,
	stationID int32,
) ([]types.Measurement, error) {
	station, err := s.store.GetStationByID(ctx, stationID)
	if err != nil {
		return nil, err
	}
	if station == nil {
		return nil, fmt.Errorf("station with ID %d not found", stationID)
	}
	measurements, err := s.store.GetLatestMeasurementsByStation(ctx, station.Id)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}
