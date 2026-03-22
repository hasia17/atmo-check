package main

import (
	"aggregator/internal/aggregator"
	"aggregator/internal/api"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	service := aggregator.NewService(ctx)

	http.HandleFunc("/aggregatedData/{voivodeship}", getAggregatedData(service))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func getAggregatedData(service *aggregator.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to get aggregated data started")
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		voivodeshipStr := r.PathValue("voivodeship")
		voivodeship, err := api.MapVoivodeship(voivodeshipStr)
		if err != nil {
			http.Error(w, "Unknown voivodeship: "+voivodeshipStr, http.StatusBadRequest)
			return
		}
		results, err := service.AggregateData(ctx, voivodeship)
		if err != nil {
			slog.Error("Aggregating data failed", "voivodeship", voivodeshipStr, "error", err)
			http.Error(w, "Aggregating data for voivodeship failed", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(results); err != nil {
			slog.Error("Encoding json response failed", "error", err)
			http.Error(w, "Encoding json response failed", http.StatusInternalServerError)
			return
		}
		slog.Info("Request to get aggregated data finished successfully")
	}
}
