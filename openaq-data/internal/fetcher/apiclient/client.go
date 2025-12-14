package apiclient

import (
	"context"
	"fmt"
	"log/slog"

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
	logger *slog.Logger
}

func New(apiKey string, l *slog.Logger) (*Service, error) {
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

func (s *Service) FetchLocations(ctx context.Context) ([]OpenAqLocation, error) {
	var data openAqLocationResponse
	queryParams := map[string]string{
		"iso":   "PL",
		"limit": "1000",
	}
	resp, err := s.request(
		ctx,
		&data,
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

	return data.Results, nil
}

func (s *Service) FetchMeasurementsForLocation(ctx context.Context, locationId, paramId int32) ([]OpenAqMeasurement, error) {
	var apiData openAqMeasurementResponse
	queryParams := map[string]string{
		"parameter_id": fmt.Sprintf("%d", paramId),
	}
	resp, err := s.request(
		ctx,
		&apiData,
		fmt.Sprintf(measurementsEndpoint, locationId),
		resty.MethodGet,
		queryParams,
	).Send()
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

func (s *Service) FetchParameters(ctx context.Context) ([]OpenAqParameter, error) {
	var data openAqParameterResponse
	res, err := s.request(
		ctx,
		&data,
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
	return data.Results, nil
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
