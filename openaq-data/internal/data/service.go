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

func (s *Service) GetStations(ctx context.Context) ([]types.Station, error) {
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		s.logger.Error("Failed to get stations", slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("Stations returned", slog.Int("count", len(stations)))
	return stations, nil
}

func (s *Service) GetStationByID(ctx context.Context, id int32) (*types.Station, error) {
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
) ([]types.Measurement, error) {
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
	return measurements, nil
}
