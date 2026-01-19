package openMeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Hostname = "http://localhost:8083"

func GetStations() ([]Station, error) {

	response, err := http.Get(Hostname + "/open-meteo-data-rs/stations")
	if err != nil {
		return nil, fmt.Errorf("request to fetch stations failed: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read stations: %v", err)
	}

	var stations []Station
	if err = json.Unmarshal(body, &stations); err != nil {
		return nil, fmt.Errorf("unmarshalling response body failed: %v", err)
	}

	return stations, nil
}

func GetParameters() ([]Parameter, error) {
	response, err := http.Get(Hostname + "/open-meteo-data-rs/parameters")
	if err != nil {
		return nil, fmt.Errorf("request to fetch parameters failed: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read parameters body: %v", err)
	}

	var parameters []Parameter
	if err = json.Unmarshal(body, &parameters); err != nil {
		return nil, fmt.Errorf("unmarshalling response body failed: %v", err)
	}
	return parameters, nil
}

func GetMeasurementForStation(stationId string) ([]Measurement, error) {
	response, err := http.Get(Hostname + "/open-meteo-data-rs/stations/" + stationId + "/measurements")
	if err != nil {
		return nil, fmt.Errorf("request to fetch measurements for station %s failed: %v", stationId, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read measurements for station %s body: %v", stationId, err)
	}

	var measurements []Measurement
	if err = json.Unmarshal(body, &measurements); err != nil {
		return nil, fmt.Errorf("unmarshalling response body for station %s failed: %v", stationId, err)
	}
	return measurements, nil
}
