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
