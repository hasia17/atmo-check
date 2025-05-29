package internal

import (
	"fmt"

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
	fmt.Println("Handling GET /data request")
	return nil
}
