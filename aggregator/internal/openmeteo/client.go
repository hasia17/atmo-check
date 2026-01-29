package openmeteo

import (
	"aggregator/internal/apiclient"
	"fmt"
)

const Hostname = "http://localhost:8083"

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetStations() ([]Station, error) {
	return apiclient.FetchData[Station](Hostname + "/open-meteo-data-rs/stations")
}

func (c *Client) GetParameters() ([]Parameter, error) {
	return apiclient.FetchData[Parameter](Hostname + "/open-meteo-data-rs/parameters")
}

func (c *Client) GetMeasurementForStation(stationId int) ([]Measurement, error) {
	url := fmt.Sprintf("%s/open-meteo-data-rs/stations/%d/measurements", Hostname, stationId)
	return apiclient.FetchData[Measurement](url)
}
