package open_meteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Hostname = "http://localhost:8083"

func GetStations() ([]Station, error) {
	return readResponse[Station](Hostname + "/open-meteo-data-rs/stations")
}

func GetParameters() ([]Parameter, error) {
	return readResponse[Parameter](Hostname + "/open-meteo-data-rs/parameters")
}

func GetMeasurementForStation(stationId string) ([]Measurement, error) {
	return readResponse[Measurement](Hostname + "/open-meteo-data-rs/stations/" + stationId + "/measurements")
}

func readResponse[T any](url string) ([]T, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %v", url, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body for request %u. Error: %v", url, err)
	}

	var results []T
	if err = json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("unmarshalling response body for request %u failed: %v", url, err)
	}
	return results, nil
}
