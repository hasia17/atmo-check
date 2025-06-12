package internal

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type DataHandler struct {
	service *DataService
	logger  *slog.Logger
}

func NewDataHandler(service *DataService, l *slog.Logger) *DataHandler {
	return &DataHandler{
		service: service,
		logger:  l,
	}
}

func (h *DataHandler) HandleGetLocations(c fiber.Ctx) error {
	h.logger.Info("Fetching locations")
	locations, err := h.service.store.GetLocations(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch locations")
	}

	h.logger.Info("Locations fetched successfully", slog.Int("count", len(locations)))
	return c.JSON(fiber.Map{
		"data": locations,
	})
}

func (h *DataHandler) HandleGetMeasurementsByLocation(c fiber.Ctx) error {
	h.logger.Info("Fetching measurements for location", slog.String("locationId", c.Params("id")))

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

	h.logger.Info(
		"Measurements fetched successfully",
		slog.Int("locationId", locationId),
		slog.Int("count", len(measurements)),
	)
	return c.JSON(fiber.Map{
		"data": measurements,
	})
}

func (h *DataHandler) HandleGetLocationByID(c fiber.Ctx) error {
	h.logger.Info("Fetching location by ID", slog.String("locationId", c.Params("id")))

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
	}
	location, err := h.service.store.GetLocationByID(c.Context(), int32(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch location")
	}
	if location == nil {
		return fiber.NewError(fiber.StatusNotFound, "Location not found")
	}

	h.logger.Info("Location fetched successfully", slog.Int("locationId", int(location.ID)))
	return c.JSON(fiber.Map{"data": location})
}
