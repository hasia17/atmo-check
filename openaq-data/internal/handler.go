package internal

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type DataHandler struct {
	service *DataService
}

func NewDataHandler(service *DataService) *DataHandler {
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
	if station == nil {
		return fiber.NewError(fiber.StatusNotFound, "Station not found")
	}

	return c.JSON(fiber.Map{"data": station})
}

// HandleGetLatestMeasurementsByStation returns latest measurements for a station
func (h *DataHandler) HandleGetLatestMeasurementsByStation(c fiber.Ctx) error {
	stationIdStr := c.Params("id")
	stationId, err := strconv.Atoi(stationIdStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}

	measurements, err := h.service.GetLatestMeasurementsByStation(c.Context(), int32(stationId))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch measurements")
	}

	return c.JSON(fiber.Map{
		"data": measurements,
	})
}

// HandleGetParametersByStation returns parameters for a station
func (h *DataHandler) HandleGetParametersByStation(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid station ID")
	}
	parameters, err := h.service.GetParametersByStationID(c.Context(), int32(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch parameters")
	}
	if parameters == nil {
		return fiber.NewError(fiber.StatusNotFound, "Station not found")
	}

	return c.JSON(fiber.Map{"data": parameters})
}
