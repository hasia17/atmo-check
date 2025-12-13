package internal

import (
	"context"
	"openaq-data/internal/types"
)

type FetcherService interface {
	Run(ctx context.Context) error
}

type DataService interface {
	Stations(ctx context.Context) ([]types.Station, error)
	Parameters(ctx context.Context) ([]types.Parameter, error)
	MeasurementsForStation(ctx context.Context, stationID int32) ([]types.Measurement, error)
}
