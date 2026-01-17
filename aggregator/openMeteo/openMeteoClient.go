package openMeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetStations() []Station {

	url := "http://localhost:8083/open-meteo-data-rs/stations"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Request to fetch stations failed", err)
		return nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Reading response body failed", err)
		return nil
	}

	var stations []Station
	err = json.Unmarshal(body, &stations)
	if err != nil {
		fmt.Println("Unmarshalling response body failed", err)
		return nil
	}

	return stations
}
