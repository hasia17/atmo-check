package main

import (
	"aggregator/openMeteo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/stations", getStations)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getStations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to get stations started")
	var stations []openMeteo.Station

	stations = openMeteo.GetStations()
	if stations == nil {
		fmt.Println("Fetching stations failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(stations)
	if err != nil {
		fmt.Println("Encoding json response failed", err)
		return
	}
	fmt.Println("Request to get stations finished successfully")

}
