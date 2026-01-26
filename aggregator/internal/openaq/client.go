package openaq

import (
	"aggregator/internal/apiclient"
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
