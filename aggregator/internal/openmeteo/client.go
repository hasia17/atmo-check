package openmeteo

import (
	"aggregator/internal/apiclient"
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

func (c *Client) GetMeasurementForStation(stationId string) ([]Measurement, error) {
	return apiclient.FetchData[Measurement](Hostname + "/open-meteo-data-rs/stations/" + stationId + "/measurements")
}
