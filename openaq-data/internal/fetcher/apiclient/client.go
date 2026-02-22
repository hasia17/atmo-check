package apiclient

import (
	"context"
	"fmt"
	"openaq-data/internal/models"

	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	baseURL = "https://api.openaq.org/v3/"
)

var (
	ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded")
)

type Service struct {
	client *resty.Client
	logger *zap.SugaredLogger
}

func New(apiKey string, l *zap.SugaredLogger) (*Service, error) {
	client, err := buildClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: client,
		logger: l,
	}, nil
}

func buildClient(apiKey string) (*resty.Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	return resty.New().
			SetHeader("Accept", "application/json").
			SetHeader("X-API-Key", apiKey),
		nil
}

func (s *Service) FetchLocations(ctx context.Context) ([]models.Location, error) {
	var locationResponse struct {
		Results []models.Location `json:"results"`
	}

	queryParams := map[string]string{
		"iso":   "PL",
		"limit": "1000",
	}
	resp, err := s.request(
		ctx,
		&locationResponse,
		locationsEndpoint,
		resty.MethodGet,
		queryParams,
	).Send()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return locationResponse.Results, nil
}

func (s *Service) FetchMeasurementsForLocation(ctx context.Context, locationId int32) ([]models.Measurement, error) {
	var measurementResponse struct {
		Results []models.Measurement `json:"results"`
	}

	resp, err := s.request(
		ctx,
		&measurementResponse,
		fmt.Sprintf(measurementsEndpoint, locationId),
		resty.MethodGet,
		nil,
	).Send()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch measurements for location %d, %w",
			locationId,
			err,
		)
	}
	if resp.IsError() {
		if resp.StatusCode() == 429 {
			return nil, ErrRateLimitExceeded
		}
		return nil, fmt.Errorf("API error for location %d: %s", locationId, resp.Status())
	}
	return measurementResponse.Results, nil
}

func (s *Service) FetchParameters(ctx context.Context) ([]models.Parameter, error) {
	var parameterResponse struct {
		Results []models.Parameter `json:"results"`
	}
	res, err := s.request(
		ctx,
		&parameterResponse,
		parametersEndpoint,
		resty.MethodGet,
		nil,
	).Send()
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	if res.IsError() {
		return nil, fmt.Errorf("API error: %s", res.Status())
	}
	return parameterResponse.Results, nil
}

func (s *Service) request(
	ctx context.Context,
	result any,
	endpoint string,
	method string,
	queryParams map[string]string,
) *resty.Request {
	return s.client.R().
		SetContext(ctx).
		SetResult(result).
		SetURL(baseURL + endpoint).
		SetMethod(method).
		SetQueryParams(queryParams)
}
