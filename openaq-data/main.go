package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"openaq-data/internal/api"
	"openaq-data/internal/fetcher"
	"openaq-data/internal/server"
	"openaq-data/internal/store"

	"github.com/joho/godotenv"
)

const (
	listenAddr = ":3000"
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

	db, err := store.New(mongoURI)
	if err != nil {
		log.Fatal("failed to open DB")
	}
	defer db.Close()

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // move this to the future config file
	}))

	f, err := fetcher.NewService(apiKey, db, l)
	if err != nil {
		l.Error("failed to create fetcher service", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	defer cancel()

	go func() {
		if err := f.Run(ctx); err != nil {
			l.Error("fetcher service stopped with error", "error", err)
			cancel()
			return
		}
	}()

	srv := server.New(db, l)
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: api.Handler(srv),
	}
	go func() {
		l.Info("Starting server", "addr", listenAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("Server failed", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		l.Error("Failed to shut down server", "error", err)
	}

	if err := db.Close(); err != nil {
		l.Error("Failed to close data service", "error", err)
	}
}
