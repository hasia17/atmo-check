package openmeteo

import (
	"aggregator/internal/apiclient"
	"context"
	"fmt"
	"os"
)

const defaultHostName = "http://localhost:8083"

type Client struct {
	hostname string
}

func NewClient() *Client {
	hostname := os.Getenv("OPENMETEO_URL")
	if hostname == "" {
		hostname = defaultHostName
	}
	return &Client{hostname: hostname}
}

func (c *Client) GetStations(ctx context.Context) ([]Station, error) {
	return apiclient.FetchData[Station](ctx, c.hostname+"/stations")
}

func (c *Client) GetParameters(ctx context.Context) ([]Parameter, error) {
	return apiclient.FetchData[Parameter](ctx, c.hostname+"/parameters")
}

func (c *Client) GetMeasurementForStation(ctx context.Context, stationId int) ([]Measurement, error) {
	url := fmt.Sprintf("%s/stations/%d/measurements", c.hostname, stationId)
	return apiclient.FetchData[Measurement](ctx, url)
}
