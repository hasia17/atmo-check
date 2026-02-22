package fetcher

import (
	"context"
	"errors"
	"fmt"
	"openaq-data/internal"
	"openaq-data/internal/fetcher/apiclient"
	"openaq-data/internal/models"
	"openaq-data/internal/store"
	"openaq-data/internal/util"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	locationsUpdateInterval     = 24 * time.Hour
	measurementsUpdateInterval  = 1 * time.Hour
	maxFetchMeasurementsRetries = 3
)

// Fetcher service is responsible for fetching data from the OpenAQ API
// and storing it in the local store.
// It runs in the background, periodically updating locations and measurements.
type Service struct {
	client *apiclient.Service
	store  store.Storer
	logger *zap.SugaredLogger

	locationsLoadedOnce  sync.Once
	locationsLoaded      chan struct{}
	parametersLoadedOnce sync.Once
	parametersLoaded     chan struct{}
}

func NewService(apiKey string, s store.Storer, l *zap.SugaredLogger) (internal.FetcherService, error) {
	apiclient, err := apiclient.New(apiKey, l)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return &Service{
		client:           apiclient,
		store:            s,
		logger:           l,
		locationsLoaded:  make(chan struct{}),
		parametersLoaded: make(chan struct{}),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	go s.updateLocationsLoop(ctx)
	go s.updateParametersLoop(ctx)
	go s.updateMeasurementsLoop(ctx)

	<-ctx.Done()
	s.logger.Info("Service is shutting down...")
	return nil
}

func (s *Service) updateLocationsLoop(ctx context.Context) {
	ticker := time.NewTicker(locationsUpdateInterval)
	defer ticker.Stop()

	for {
		if err := s.loadLocations(ctx); err != nil {
			s.logger.Errorw("Failed to update locations", "error", err)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) loadLocations(ctx context.Context) error {
	locations, err := s.client.FetchLocations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch locations from API: %w", err)
	}
	if err := s.store.StoreLocations(ctx, locations); err != nil {
		return fmt.Errorf("failed to store locations: %w", err)
	}

	s.locationsLoadedOnce.Do(func() {
		close(s.locationsLoaded)
	})

	return nil
}

func (s *Service) updateParametersLoop(ctx context.Context) {
	ticker := time.NewTicker(locationsUpdateInterval)
	defer ticker.Stop()

	for {
		if err := s.loadParameters(ctx); err != nil {
			s.logger.Errorw("Failed to update parameters", "error", err)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) loadParameters(ctx context.Context) error {
	parameters, err := s.client.FetchParameters(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch parameters from API: %w", err)
	}
	if err := s.store.StoreParameters(ctx, parameters); err != nil {
		return fmt.Errorf("failed to store parameters: %w", err)
	}

	s.parametersLoadedOnce.Do(func() {
		close(s.parametersLoaded)
	})

	return nil
}

func (s *Service) updateMeasurementsLoop(ctx context.Context) {
	select {
	case <-util.WaitFor(s.locationsLoaded, s.parametersLoaded):
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(measurementsUpdateInterval)
	defer ticker.Stop()
	for {
		if err := s.updateMeasurements(ctx); err != nil {
			s.logger.Errorw("Failed to update measurements", "error", err)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateMeasurements(ctx context.Context) error {
	locations, err := s.store.GetLocations(ctx)
	if err != nil {
		return err
	}
	if len(locations) == 0 {
		return fmt.Errorf("no locations found in store")
	}
	for _, loc := range locations {
		if err := s.updateMeasurementsForLocation(ctx, loc); err != nil {
			s.logger.Errorw(
				"Failed to update measurements for location",
				"error", err,
			)
		}
	}
	return nil
}

func (s *Service) updateMeasurementsForLocation(ctx context.Context, loc models.Location) error {
	measurements, err := s.loadMeasurementsForLocation(ctx, loc.Id)
	if err != nil {
		return err
	}
	if len(measurements) == 0 {
		return fmt.Errorf("no measurements found for location %s", loc.Name)
	}

	err = s.store.DeleteMeasurementsForLocation(ctx, loc.Id)
	if err != nil {
		return fmt.Errorf("failed to delete existing measurements for location %s: %w", loc.Name, err)
	}

	err = s.store.StoreMeasurements(ctx, measurements)
	if err != nil {
		return fmt.Errorf("failed to store measurements for location %s: %w", loc.Name, err)
	}
	return nil
}

func (s *Service) loadMeasurementsForLocation(ctx context.Context, locId int32) ([]models.Measurement, error) {
	for range maxFetchMeasurementsRetries {
		apiData, err := s.client.FetchMeasurementsForLocation(ctx, locId)
		if err != nil {
			if errors.Is(err, apiclient.ErrRateLimitExceeded) {
				s.logger.Warn("Rate limit exceeded, retrying after delay")
				<-time.After(5 * time.Second)
				continue
			}
			return nil, err
		}
		return apiData, nil
	}
	return nil, fmt.Errorf("failed to fetch measurements for location %d after retries", locId)
}
