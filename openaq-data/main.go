package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"openaq-data/internal"
	"openaq-data/internal/data"

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

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // move this to the future config file
	}))
	s, err := data.NewService(apiKey, mongoURI, l)
	if err != nil {
		l.Error("Failed to create data service", "error", err)
		os.Exit(1)
	}
	defer s.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := s.Run(ctx); err != nil {
			l.Error("Data service stopped with error", "error", err)
			cancel()
			return
		}
	}()

	h := internal.NewDataHandler(s)
	app := fiber.New()
	app.Get("/stations", h.HandleGetStations)
	app.Get("/stations/:id", h.HandleGetStationByID)
	app.Get("/stations/:id/measurements", h.HandleGetMeasurementsForStation)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			l.Error("Failed to start server", "error", err)
			cancel()
			return
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	if err := app.Shutdown(); err != nil {
		l.Error("Failed to shutdown server", "error", err)
	}

	if err := s.Close(); err != nil {
		l.Error("Failed to close data service", "error", err)
	}
}
