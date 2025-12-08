package api

import (
	"context"
	"fmt"
	"log/slog"

	"resty.dev/v3"
)

var (
	ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded")
)

type Service struct {
	apiKey string
	client *resty.Client
	logger *slog.Logger
}

func New(apiKey string, l *slog.Logger) *Service {
	client := resty.New().
		SetBaseURL("https://api.openaq.org/v3/").
		SetHeader("Content-Type", "application/json")
	if apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}
	return &Service{
		apiKey: apiKey,
		client: client,
		logger: l,
	}
}

func (s *Service) FetchLocations(ctx context.Context) ([]OpenAQLocation, error) {
	var data openAQLocationResponse
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"iso":   "PL",
			"limit": "1000",
		}).
		SetResult(&data).
		Get("locations")

	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return data.Results, nil
}

func (s *Service) FetchMeasurementsForLocation(locationId, paramId int32) ([]OpenAQMeasurement, error) {
	var apiData openAQMeasurementResponse
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"parameter_id": fmt.Sprintf("%d", paramId),
		}).
		SetResult(&apiData).
		Get(fmt.Sprintf("locations/%d/latest", locationId))

	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch measurements for location %d, parameter %d: %w",
			locationId,
			paramId,
			err,
		)
	}
	if resp.IsError() {
		if resp.StatusCode() == 429 {
			return nil, ErrRateLimitExceeded
		}
		return nil, fmt.Errorf("API error for location %d, parameter %d: %s", locationId, paramId, resp.Status())
	}
	return apiData.Results, nil
}
