package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"openaq-data/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it.")
	}

	apiKey := os.Getenv("OPENAQ_API_KEY")
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is required")
	}

	s, err := internal.NewDataService(apiKey, mongoURI)
	if err != nil {
		log.Fatalf("Failed to create data service: %v", err)
	}
	defer s.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := s.Run(ctx); err != nil {
			log.Printf("Data service run failed: %v", err)
		}
	}()

	h := internal.NewDataHandler(s)
	app := fiber.New()
	app.Get("/data", func(c fiber.Ctx) error {
		return h.HandleGetData(c)
	})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	if err := s.Close(); err != nil {
		log.Printf("Data service close failed: %v", err)
	}
}
