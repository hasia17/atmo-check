package mock

import (
	"context"
	"openaq-data/internal/models"
	"openaq-data/internal/store"
)

type Store struct {
	Locations    []models.Location
	Measurements []models.Measurement
	Parameters   []models.Parameter
}

func New() store.Storer {
	return &Store{
		Locations:    []models.Location{},
		Measurements: []models.Measurement{},
		Parameters:   []models.Parameter{},
	}
}

func (s *Store) StoreLocations(_ context.Context, locations []models.Location) error {
	s.Locations = locations
	return nil
}

func (s *Store) GetLocations(_ context.Context) ([]models.Location, error) {
	return s.Locations, nil
}

func (s *Store) GetLocationByID(_ context.Context, id int32) (*models.Location, error) {
	for _, loc := range s.Locations {
		if loc.Id == id {
			return &loc, nil
		}
	}
	return nil, nil
}

func (s *Store) StoreMeasurements(_ context.Context, measurements []models.Measurement) error {
	s.Measurements = measurements
	return nil
}

func (s *Store) DeleteMeasurementsForLocation(_ context.Context, locId int32) error {
	filtered := s.Measurements[:0] // keep the same underlying array
	for _, m := range s.Measurements {
		if m.LocationId != locId {
			filtered = append(filtered, m)
		}
	}
	s.Measurements = filtered
	return nil
}

func (s *Store) GetMeasurementsByLocation(_ context.Context, locationId int32) ([]models.Measurement, error) {
	var result []models.Measurement
	for _, m := range s.Measurements {
		if m.LocationId == locationId {
			result = append(result, m)
		}
	}
	return result, nil
}

func (s *Store) StoreParameters(_ context.Context, parameters []models.Parameter) error {
	s.Parameters = parameters
	return nil
}

func (s *Store) GetParameters(_ context.Context) ([]models.Parameter, error) {
	return s.Parameters, nil
}

func (s *Store) Close() error {
	return nil
}
