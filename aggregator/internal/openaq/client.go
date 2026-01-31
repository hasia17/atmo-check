package openaq

import (
	"aggregator/internal/apiclient"
	"fmt"
)

const Hostname = "http://localhost:3000"

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetStations() ([]Station, error) {

	return apiclient.FetchData[Station](Hostname + "/stations")
}

func (c *Client) GetMeasurementForStation(stationId int) ([]Measurement, error) {
	url := fmt.Sprintf("%s/stations/%d/measurements", Hostname, stationId)
	return apiclient.FetchData[Measurement](url)
}
