package server

import (
	"log/slog"
	"openaq-data/internal/data"
	"openaq-data/internal/store"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

const (
	listenAddr = ":3000"
)

type Server struct {
	app         *fiber.App
	dataService *data.Service
}

func New(db *store.Store, l *slog.Logger) *Server {
	d := data.NewService(db, l)

	s := &Server{
		dataService: d,
	}

	app := fiber.New()
	app.Get("/stations", s.handleGetStations)
	app.Get("/stations/:id", s.handleGetStationByID)
	app.Get("/stations/:id/measurements", s.handleGetMeasurementsForStation)
	s.app = app

	return s
}

func (s *Server) Run() error {
	return s.app.Listen(listenAddr)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

func (s *Server) handleGetStations(c fiber.Ctx) error {
	stations, err := s.dataService.Stations(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch stations")
	}

	return c.JSON(fiber.Map{
		"data": stations,
	})
}

func (s *Server) handleGetStationByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}
	station, err := s.dataService.StationByID(c.Context(), int32(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch station")
	}
	return c.JSON(fiber.Map{"data": station})
}

func (s *Server) handleGetMeasurementsForStation(c fiber.Ctx) error {
	stationIdStr := c.Params("id")
	stationId, err := strconv.Atoi(stationIdStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	measurements, err := s.dataService.MeasurementsForStation(c.Context(), int32(stationId), 100)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch measurements")
	}

	return c.JSON(fiber.Map{
		"data": measurements,
	})
}
