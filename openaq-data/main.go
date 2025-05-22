package main

import (
	"log"
	"openaq-data/internal"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it.")
	}

	apiKey := ""
	if v := os.Getenv("OPENAQ_API_KEY"); v != "" {
		apiKey = v
	}
	ls := internal.NewLocationService()
	s := internal.NewDataService(apiKey, ls)
	h := internal.NewDataHandler(s)

	app := fiber.New()
	app.Get("/data", func(c fiber.Ctx) error {
		return h.HandleGetData(c, c.Query("location"))
	})
	log.Fatal(app.Listen(":3000"))
}
