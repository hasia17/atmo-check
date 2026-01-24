package openaq

import "aggregator/internal/util"

const Hostname = "http://localhost:3000"

func GetStations() ([]Station, error) {
	return util.ReadResponse[Station](Hostname + "/stations")
}
