package internal

import (
	"context"
	"openaq-data/internal/api"
)

type FetcherService interface {
	Run(ctx context.Context) error
}

type DataService interface {
	Stations(ctx context.Context) ([]api.Station, error)
	Parameters(ctx context.Context) ([]api.Parameter, error)
	MeasurementsForStation(ctx context.Context, stationID int32) ([]api.Measurement, error)
}
