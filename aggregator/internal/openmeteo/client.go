package openmeteo

import (
	"aggregator/internal/apiclient"
	"context"
	"fmt"
)

const Hostname = "http://localhost:8083"

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetStations(ctx context.Context) ([]Station, error) {
	return apiclient.FetchData[Station](ctx, Hostname+"/stations")
}

func (c *Client) GetParameters(ctx context.Context) ([]Parameter, error) {
	return apiclient.FetchData[Parameter](ctx, Hostname+"/parameters")
}

func (c *Client) GetMeasurementForStation(ctx context.Context, stationId int) ([]Measurement, error) {
	url := fmt.Sprintf("%s/stations/%d/measurements", Hostname, stationId)
	return apiclient.FetchData[Measurement](ctx, url)
}
