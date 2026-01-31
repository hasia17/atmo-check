package main

import (
	"aggregator/internal/aggregator"
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	opmClient := openmeteo.NewClient()
	oaqClient := openaq.NewClient()

	service, err := aggregator.NewService(opmClient, oaqClient)
	if err != nil {
		log.Fatal(err)
	}
	results, err := service.AggregateOpenMeteo(api.Dolnoslaskie)
	log.Println(results)

	http.HandleFunc("/stations", getStations)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getStations(w http.ResponseWriter, r *http.Request) {
	log.Print("Request to get stations started")
	var stations []openmeteo.Station

	openmeteoClient := openmeteo.NewClient()

	stations, err := openmeteoClient.GetStations()
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
