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
// @ID getStations
// @Summary Get all stations
// @Description Returns a list of all stations
// @Tags stations
// @Produce json
// @Success 200 {object} map[string][]Station
// @Router /stations [get]
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
// @ID getStationById
// @Summary Get station by ID
// @Description Returns a station by its ID
// @Tags stations
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} map[string]Station
// @Failure 404 {object} map[string]string
// @Router /stations/{id} [get]
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

// HandleGetMeasurementsByStation returns latest measurements for a station
// @ID getMeasurementsByStation
// @Summary Get latest measurements by station
// @Description Returns the latest measurement for each parameter at a specific station
// @Tags measurements
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} map[string][]Measurement
// @Router /stations/{id}/measurements [get]
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
// @ID getParametersByStation
// @Summary Get parameters by station
// @Description Returns parameters for a specific station
// @Tags parameters
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} map[string][]Parameter
// @Router /stations/{id}/parameters [get]
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
