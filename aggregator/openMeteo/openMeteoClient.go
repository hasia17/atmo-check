package openMeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const HOST_NAME = "http://localhost:8083"

func GetStations() ([]Station, error) {

	url := "/open-meteo-data-rs/stations"

	response, err := http.Get(HOST_NAME + url)
	if err != nil {
		return nil, fmt.Errorf("request to fetch stations failed: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read stations: %v", err)
	}

	var stations []Station
	err = json.Unmarshal(body, &stations)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body failed: %v", err)
	}

	return stations, nil
}
