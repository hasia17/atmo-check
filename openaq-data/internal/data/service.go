package data

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"openaq-data/internal"
	"openaq-data/internal/api"
	"openaq-data/internal/store"
	"openaq-data/internal/types"
	"openaq-data/internal/util"
)

type Service struct {
	store  store.Storer
	logger *slog.Logger
}

func NewService(s store.Storer, l *slog.Logger) internal.DataService {
	return &Service{
		store:  s,
		logger: l,
	}
}

func (s *Service) Stations(ctx context.Context) ([]api.Station, error) {
	locations, err := s.store.GetLocations(ctx)
	if err != nil {
		return nil, err
	}
	return s.buildStations(locations), nil
}

func (s *Service) buildStations(locations []types.Location) []api.Station {
	var stations []api.Station
	for _, apiLoc := range locations {
		station := api.Station{
			Id:        apiLoc.Id,
			Name:      apiLoc.Name,
			Locality:  apiLoc.Locality,
			Timezone:  apiLoc.Timezone,
			Latitude:  apiLoc.Coordinates.Latitude,
			Longitude: apiLoc.Coordinates.Longitude,
		}
		for _, apiSensor := range apiLoc.Sensors {
			station.ParameterIds = append(station.ParameterIds, apiSensor.Parameter.Id)
		}
		station.ParameterIds = util.RemoveDuplicates(station.ParameterIds)
		stations = append(stations, station)
	}
	return stations
}

func (s *Service) Parameters(ctx context.Context) ([]api.Parameter, error) {
	parameters, err := s.store.GetParameters(ctx)
	if err != nil {
		return nil, err
	}
	return s.buildParameters(parameters), nil
}

func (s *Service) buildParameters(params []types.Parameter) []api.Parameter {
	var apiParams []api.Parameter
	for _, param := range params {
		apiParam := api.Parameter{
			Id:          param.Id,
			Name:        param.Name,
			Units:       param.Units,
			DisplayName: param.DisplayName,
			Description: &param.Description,
		}
		apiParams = append(apiParams, apiParam)
	}
	return apiParams
}

func (s *Service) MeasurementsForStation(
	ctx context.Context,
	locId int32,
) ([]api.Measurement, error) {
	loc, err := s.store.GetLocationByID(ctx, locId)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, fmt.Errorf("location with ID %d not found", locId)
	}
	measurements, err := s.store.GetMeasurementsByLocation(ctx, locId)
	if err != nil {
		return nil, err
	}
	return s.buildMeasurements(measurements, *loc), nil
}

func (s *Service) buildMeasurements(
	measurements []types.Measurement,
	loc types.Location,
) []api.Measurement {
	var apiMeasurements []api.Measurement
	for _, m := range measurements {
		parsedTime, err := util.StringToTime(m.Date.Utc)
		if err != nil {
			log.Printf("Failed to parse time for measurement (UTC: %s): %v", m.Date.Utc, err)
			continue
		}

		paramId, err := exctractParameterId(m.SensorId, loc)
		if err != nil {
			log.Printf("Failed to extract parameter ID for sensor ID %d: %v", m.SensorId, err)
			continue
		}

		apiMeasurements = append(apiMeasurements, api.Measurement{
			Timestamp:   parsedTime,
			Value:       m.Value,
			ParameterId: paramId,
			StationId:   loc.Id,
		})
	}
	return apiMeasurements
}

func exctractParameterId(sensorId int32, loc types.Location) (int32, error) {
	for _, sensor := range loc.Sensors {
		if sensor.Id == sensorId {
			return sensor.Parameter.Id, nil
		}
	}
	return 0, fmt.Errorf("parameter ID not found for sensor ID %d", sensorId)
}
