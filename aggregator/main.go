package main

import (
	open_meteo2 "aggregator/internal/open_meteo"
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/stations", getStations)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getStations(w http.ResponseWriter, r *http.Request) {
	log.Print("Request to get stations started")
	var stations []open_meteo2.Station

	stations, err := open_meteo2.GetStations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(stations); err != nil {
		http.Error(w, "Encoding json response failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Print("Request to get stations finished successfully")

}
