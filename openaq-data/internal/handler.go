package internal

import (
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

func (h *DataHandler) HandleGetData(c fiber.Ctx) error {
	h.service.FetchLocations()
	h.service.FetchMeasurements()
	var retData []measurement
	for _, loc := range h.service.locations {
		data := h.service.measurements[loc]
		if len(data) > 0 {
			for _, m := range data {
				retData = append(retData, m)
			}
		}
	}
	return c.JSON(retData)
}
