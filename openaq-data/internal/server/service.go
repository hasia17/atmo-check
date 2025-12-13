package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"openaq-data/internal"
	"openaq-data/internal/api"
)

type Service struct {
	dataService internal.DataService
	logger      *slog.Logger
}

func New(d internal.DataService, l *slog.Logger) api.ServerInterface {
	return &Service{
		dataService: d,
		logger:      l,
	}
}

func (s *Service) GetStations(w http.ResponseWriter, r *http.Request) {
	stations, err := s.dataService.Stations(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch stations", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data": stations,
	})
}

func (s *Service) GetParameters(w http.ResponseWriter, r *http.Request) {
	parameters, err := s.dataService.Parameters(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch parameters", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusNotImplemented, map[string]any{
		"data": parameters,
	})
}

func (s *Service) GetMeasurementsByStation(w http.ResponseWriter, r *http.Request, id int32) {
	measurements, err := s.dataService.MeasurementsForStation(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch measurements", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data": measurements,
	})
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
