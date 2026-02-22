package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"openaq-data/internal/api"
	"openaq-data/internal/data"
	"openaq-data/internal/fetcher"
	"openaq-data/internal/server"
	"openaq-data/internal/store"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	listenAddr = ":3000"
)

func makeLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger.Sugar()
}

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

	l := makeLogger()

	f, err := fetcher.NewService(apiKey, db, l)
	if err != nil {
		l.Errorw("failed to create fetcher service", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	defer cancel()

	go func() {
		if err := f.Run(ctx); err != nil {
			l.Errorw("fetcher service stopped with error", "error", err)
			cancel()
			return
		}
	}()

	dataService := data.NewService(db, l)
	srv := server.New(dataService, l)
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: api.Handler(srv),
	}
	go func() {
		l.Infow("Starting server", "addr", listenAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Errorw("Server failed", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	l.Info("Shutting down...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		l.Errorw("Failed to shut down server", "error", err)
	}

	if err := db.Close(); err != nil {
		l.Errorw("Failed to close data service", "error", err)
	}
}
