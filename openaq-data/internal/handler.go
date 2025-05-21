package internal

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

type DataHandler struct {
	service *Service
}

func NewDataHandler(service *Service) *DataHandler {
	return &DataHandler{
		service: service,
	}
}

func (h *DataHandler) HandleGetData(c fiber.Ctx, location string) error {
	data, err := h.service.FetchData(location)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error fetching data")
	}
	return c.JSON(data)
}
