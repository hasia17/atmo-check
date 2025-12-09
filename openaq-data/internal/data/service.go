package data

import (
	"context"
	"fmt"
	"log/slog"
	"openaq-data/internal/store"
	"openaq-data/types"
)

type Service struct {
	store  *store.Store
	logger *slog.Logger
}

func NewService(s *store.Store, l *slog.Logger) *Service {
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

func (s *Service) StationByID(ctx context.Context, id int32) (*types.Station, error) {
	station, err := s.store.GetStationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (s *Service) MeasurementsForStation(
	ctx context.Context,
	stationID int32,
	limit int64,
) ([]types.Measurement, error) {
	station, err := s.store.GetStationByID(ctx, stationID)
	if err != nil {
		return nil, err
	}
	if station == nil {
		return nil, fmt.Errorf("station with ID %d not found", stationID)
	}
	measurements, err := s.store.GetLatestMeasurementsByStation(ctx, *station.Id)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}
