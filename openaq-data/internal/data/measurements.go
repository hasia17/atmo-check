package data

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"openaq-data/api"
	"time"
)

func (s *Service) updateMeasurementsLoop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			s.logger.Info("Measurements update loop stopped")
			return
		}

		s.updateMeasurements(ctx)

		timer := time.NewTimer(time.Hour)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			s.logger.Info("Measurements update loop stopped")
			return
		case <-timer.C:
		}
	}
}

func (s *Service) updateMeasurements(ctx context.Context) {
	s.logger.Info("Updating measurements for all stations...")
	stations, err := s.store.GetStations(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to fetch stations from store: %v", err))
		return
	}
	if len(stations) == 0 {
		s.logger.Info("No stations found in store, skipping measurements fetch")
		return
	}
	s.logger.Info("Fetching locations and measurements...")
	for _, station := range stations {
		s.updateMeasurementsForStation(ctx, station)
	}
}

func (s *Service) updateMeasurementsForStation(ctx context.Context, station api.Station) {
	measurements, err := s.fetchMeasurementsForStation(station)
	if err != nil {
		s.logger.Error(
			"Failed to fetch measurements for station",
			slog.String("station", *station.Name),
			slog.Any("error", err),
		)
		return
	}
	if len(measurements) == 0 {
		s.logger.Info("No measurements found for station", slog.String("station", *station.Name))
		return
	}

	err = s.store.DeleteMeasurementsForStation(ctx, *station.Id)
	if err != nil {
		s.logger.Error(
			"Failed to delete old measurements for station",
			slog.String("station", *station.Name),
			slog.Any("error", err),
		)
		return
	}

	err = s.store.StoreMeasurements(ctx, measurements)
	if err != nil {
		s.logger.Error(
			"Failed to upsert measurements for station",
			slog.String("station", *station.Name),
			slog.Any("error", err),
		)
		return
	}
	s.logger.Info(
		"Updated measurements for station",
		slog.String("station", *station.Name),
		slog.Int("count", len(measurements)),
	)
}

func (s *Service) fetchMeasurementsForStation(st api.Station) ([]api.Measurement, error) {
	s.logger.Info("Updating measurements for station", slog.String("station", *st.Name))

	measurements := make([]api.Measurement, 0)
	for _, parameter := range *st.Parameters {
		apiData, err := s.try(*parameter.Id, *st.Id)
		if err != nil {
			s.logger.Error(
				"Failed to fetch measurements for station with parameter",
				slog.String("station", *st.Name),
				slog.String("parameter", *parameter.Name),
				slog.Any("error", err),
			)
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
			measurements = append(measurements, measurement)
		}
		<-time.After(1 * time.Second)
	}
	return measurements, nil
}

func (s *Service) try(paramId, stationId int32) (openAQMeasurementResponse, error) {
	var apiData openAQMeasurementResponse
	for range 3 {
		resp, err := s.client.R().
			SetQueryParams(map[string]string{
				"parameter_id": fmt.Sprintf("%d", paramId),
			}).
			SetResult(&apiData).
			Get(fmt.Sprintf("locations/%d/latest", stationId))

		if err != nil {
			return apiData, fmt.Errorf(
				"failed to fetch measurements for station %d, parameter %d: %w",
				stationId,
				paramId,
				err,
			)
		}
		if resp.IsError() {
			if resp.StatusCode() == 429 {
				log.Printf("Rate limit exceeded, retrying...")
				<-time.After(time.Duration(1+time.Now().UnixNano()%5) * time.Second)
				continue
			}
			return apiData, fmt.Errorf("API error for station %d, parameter %d: %s", stationId, paramId, resp.Status())
		}
	}
	return apiData, nil
}
