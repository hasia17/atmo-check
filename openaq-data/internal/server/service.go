package server

import (
	"log/slog"
	"net/http"
	"openaq-data/internal/api"
	"openaq-data/internal/data"
	"openaq-data/internal/store"

	"github.com/labstack/echo/v4"
)

type Service struct {
	dataService *data.Service
	logger      *slog.Logger
}

func New(db *store.Store, l *slog.Logger) api.ServerInterface {
	return &Service{
		dataService: data.NewService(db, l),
	}
}

func (s *Service) GetStations(ctx echo.Context) error {
	stations, err := s.dataService.Stations(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stations")
	}

	return ctx.JSON(200, echo.Map{
		"data": stations,
	})
}

func (s *Service) GetParameters(ctx echo.Context) error {
	return ctx.JSON(400, echo.Map{"message": "Not implemented"})
}

func (s *Service) GetMeasurementsByStation(ctx echo.Context, id int32) error {
	measurements, err := s.dataService.MeasurementsForStation(ctx.Request().Context(), id, 100)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch measurements")
	}

	return ctx.JSON(200, echo.Map{
		"data": measurements,
	})
}
