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

func (h *DataHandler) HandleGetLocations(c fiber.Ctx) error {
	locations, err := h.service.store.GetLocations(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch locations")
	}
	return c.JSON(fiber.Map{
		"data": locations,
	})
}

func (h *DataHandler) HandleGetMeasurementsByLocation(c fiber.Ctx) error {
	locationIdStr := c.Params("id")
	locationId, err := strconv.Atoi(locationIdStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
	}

	limitStr := c.Query("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit parameter")
	}
	if limit < 1 {
		limit = 100
	}

	measurements, err := h.service.store.GetMeasurementsByLocation(c.Context(), int32(locationId), int64(limit))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch measurements")
	}
	return c.JSON(fiber.Map{
		"data": measurements,
	})
}
