package main

import (
	"aggregator/internal/aggregator"
	"aggregator/internal/api"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	service, err := aggregator.NewService(ctx)
	if err != nil {
		log.Fatal("Failed to initialize aggregator service: ", err)
	}

	http.HandleFunc("/aggregatedData/", getAggregatedData(service))
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func getAggregatedData(service *aggregator.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request to get aggregated data started")
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		voivodeshipStr := strings.TrimPrefix(r.URL.Path, "/aggregatedData/")
		if voivodeshipStr == "" {
			http.Error(w, "Voivodeship parameter is required", http.StatusBadRequest)
			return
		}
		voivodeship, err := api.MapVoivodeship(voivodeshipStr)
		if err != nil {
			http.Error(w, "Unknown voivodeship: "+voivodeshipStr, http.StatusBadRequest)
			return
		}
		results, err := service.AggregateData(ctx, voivodeship)
		if err != nil {
			http.Error(w, "Aggregating data for voivodeship failed", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Encoding json response failed", http.StatusInternalServerError)
			return
		}
		log.Print("Request to get aggregated data finished successfully")
	}
}
