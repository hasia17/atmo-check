package internal

import (
	"openaq-data/internal/data"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type DataHandler struct {
	service *data.Service
}

func NewDataHandler(service *data.Service) *DataHandler {
	return &DataHandler{
		service: service,
	}
}

// HandleGetStations returns a list of all stations
func (h *DataHandler) HandleGetStations(c fiber.Ctx) error {
	stations, err := h.service.GetStations(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch stations")
	}

	return c.JSON(fiber.Map{
		"data": stations,
	})
}

// HandleGetStationByID returns a station by its ID
func (h *DataHandler) HandleGetStationByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}
	station, err := h.service.GetStationByID(c.Context(), int32(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch station")
	}
	return c.JSON(fiber.Map{"data": station})
}

// HandleGetLatestMeasurementsByStation returns latest measurements for a station
func (h *DataHandler) HandleGetMeasurementsForStation(c fiber.Ctx) error {
	stationIdStr := c.Params("id")
	stationId, err := strconv.Atoi(stationIdStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	measurements, err := h.service.GetMeasurementsForStation(c.Context(), int32(stationId), 100)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch measurements")
	}

	return c.JSON(fiber.Map{
		"data": measurements,
	})
}
