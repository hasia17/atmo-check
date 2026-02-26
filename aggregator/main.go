package main

import (
	"aggregator/internal/aggregator"
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/aggregate/", getAggregatedData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAggregatedData(w http.ResponseWriter, r *http.Request) {
	log.Print("Request to get aggregated data started")
	voivodeshipStr := strings.TrimPrefix(r.URL.Path, "/aggregate/")
	if voivodeshipStr == "" {
		http.Error(w, "Voivodeship parameter is required", http.StatusBadRequest)
		return
	}
	voivodeship, err := api.MapVoivodeship(voivodeshipStr)
	if err != nil {
		http.Error(w, "Unknown voivodeship: "+voivodeshipStr, http.StatusBadRequest)
		return
	}
	opmClient := openmeteo.NewClient()
	oaqClient := openaq.NewClient()
	service, err := aggregator.NewService(opmClient, oaqClient)
	if err != nil {
		http.Error(w, "Aggregator initialization failed", http.StatusInternalServerError)
		return
	}
	results, err := service.AggregateData(voivodeship)
	if err != nil {
		http.Error(w, "Aggregating data for voivodeship failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Encoding json response failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Print("Request to get aggregated data finished successfully")
}
