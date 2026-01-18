package main

import (
	"aggregator/openMeteo"
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
	var stations []openMeteo.Station

	stations, err := openMeteo.GetStations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(stations)
	if err != nil {
		http.Error(w, "Encoding json response failed", http.StatusInternalServerError)
		return
	}
	log.Print("Request to get stations finished successfully")

}
