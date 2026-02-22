package server

import (
	"encoding/json"
	"net/http"
	"openaq-data/internal"
	"openaq-data/internal/api"

	"go.uber.org/zap"
)

type Service struct {
	dataService internal.DataService
	logger      *zap.SugaredLogger
}

func New(d internal.DataService, l *zap.SugaredLogger) api.ServerInterface {
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
	writeJSON(w, http.StatusOK, stations)
}

func (s *Service) GetParameters(w http.ResponseWriter, r *http.Request) {
	parameters, err := s.dataService.Parameters(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch parameters", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, parameters)
}

func (s *Service) GetMeasurementsByStation(w http.ResponseWriter, r *http.Request, id int32) {
	measurements, err := s.dataService.MeasurementsForStation(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch measurements", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, measurements)
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(v)
}
