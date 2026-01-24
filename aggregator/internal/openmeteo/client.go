package openmeteo

import (
	"aggregator/internal/util"
)

const Hostname = "http://localhost:8083"

func GetStations() ([]Station, error) {
	return util.ReadResponse[Station](Hostname + "/open-meteo-data-rs/stations")
}

func GetParameters() ([]Parameter, error) {
	return util.ReadResponse[Parameter](Hostname + "/open-meteo-data-rs/parameters")
}

func GetMeasurementForStation(stationId string) ([]Measurement, error) {
	return util.ReadResponse[Measurement](Hostname + "/open-meteo-data-rs/stations/" + stationId + "/measurements")
}
