package aggregator

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type cache struct {
	openMeteoMap        Map[openmeteo.Station]
	openaqMap           Map[openaq.Station]
	openMeteoParameters []openmeteo.Parameter
	openaqParameters    []openaq.Parameter
	err                 error
}

const cacheRefreshInterval = 24 * time.Hour

type Service struct {
	openmeteoClient   *openmeteo.Client
	openaqClient      *openaq.Client
	voivodeshipBounds map[api.Voivodeship]geographicalBounds
	mu                sync.RWMutex
	cache             cache
}

func NewService(ctx context.Context) *Service {
	s := &Service{
		openmeteoClient: openmeteo.NewClient(),
		openaqClient:    openaq.NewClient(),
	}
	bounds, err := loadVoivodeshipBounds("config/voivodeships.json")
	if err != nil {
		s.updateCacheErr(fmt.Errorf("failed to load voivodeship bounds: %w", err))
	} else {
		s.voivodeshipBounds = bounds
	}
	go s.refreshCacheInLoop(ctx)
	return s
}

func (s *Service) refreshCacheInLoop(ctx context.Context) {
	ticker := time.NewTicker(cacheRefreshInterval)
	defer ticker.Stop()
	for {
		if err := s.refreshCache(ctx); err != nil {
			slog.Error("Failed to refresh cache", "error", err)
			s.updateCacheErr(err)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) refreshCache(ctx context.Context) error {
	var (
		openMeteoMap        Map[openmeteo.Station]
		openaqMap           Map[openaq.Station]
		openMeteoParameters []openmeteo.Parameter
		openaqParameters    []openaq.Parameter
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		stations, err := s.openmeteoClient.GetStations(gctx)
		if err != nil {
			return fmt.Errorf("fetching openmeteo stations: %w", err)
		}
		openMeteoMap = groupStationsByVoivodeship(stations, s.voivodeshipBounds)
		return nil
	})

	g.Go(func() error {
		stations, err := s.openaqClient.GetStations(gctx)
		if err != nil {
			return fmt.Errorf("fetching openaq stations: %w", err)
		}
		openaqMap = groupStationsByVoivodeship(stations, s.voivodeshipBounds)
		return nil
	})

	g.Go(func() error {
		params, err := s.openmeteoClient.GetParameters(gctx)
		if err != nil {
			return fmt.Errorf("fetching openmeteo parameters: %w", err)
		}
		openMeteoParameters = params
		return nil
	})

	g.Go(func() error {
		params, err := s.openaqClient.GetParameters(gctx)
		if err != nil {
			return fmt.Errorf("fetching openaq parameters: %w", err)
		}
		openaqParameters = params
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	s.updateCache(cache{
		openMeteoMap:        openMeteoMap,
		openaqMap:           openaqMap,
		openMeteoParameters: openMeteoParameters,
		openaqParameters:    openaqParameters,
	})
	return nil
}

func (s *Service) readCache() cache {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cache
}

func (s *Service) updateCache(c cache) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = c
}

func (s *Service) updateCacheErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache.err = err
}

type Map[T any] map[api.Voivodeship][]T

type geographicalBounds struct {
	MaxLatitude  float64 `json:"maxLat"`
	MinLatitude  float64 `json:"minLat"`
	MaxLongitude float64 `json:"maxLon"`
	MinLongitude float64 `json:"minLon"`
}

func loadVoivodeshipBounds(path string) (map[api.Voivodeship]geographicalBounds, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading voivodeship bounds file: %w", err)
	}
	var raw map[string]geographicalBounds
	if err = json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing voivodeship bounds file: %w", err)
	}
	bounds := make(map[api.Voivodeship]geographicalBounds, len(raw))
	for k, v := range raw {
		bounds[api.Voivodeship(k)] = v
	}
	return bounds, nil
}

type locatable interface {
	Latitude() float64
	Longitude() float64
	StationName() string
}

func groupStationsByVoivodeship[T locatable](stations []T, bounds map[api.Voivodeship]geographicalBounds) Map[T] {
	vm := make(map[api.Voivodeship][]T)
	for v, b := range bounds {
		for _, s := range stations {
			if stationInVoivodeship(s, b) {
				vm[v] = append(vm[v], s)
			}
		}
	}
	return vm
}

func stationInVoivodeship[T locatable](s T, b geographicalBounds) bool {
	return s.Latitude() >= b.MinLatitude &&
		s.Latitude() <= b.MaxLatitude &&
		s.Longitude() >= b.MinLongitude &&
		s.Longitude() <= b.MaxLongitude
}

func (s *Service) AggregateAll(ctx context.Context) ([]api.AggregatedData, error) {
	voivodeships := []api.Voivodeship{
		api.Dolnoslaskie, api.KujawskoPomorskie, api.Lubelskie, api.Lubuskie,
		api.Lodzkie, api.Malopolskie, api.Mazowieckie, api.Opolskie,
		api.Podkarpackie, api.Podlaskie, api.Pomorskie, api.Slaskie,
		api.Swietokrzyskie, api.WarminskoMazurskie, api.Wielkopolskie, api.Zachodniopomorskie,
	}

	results := make([]api.AggregatedData, len(voivodeships))
	g, ctx := errgroup.WithContext(ctx)
	for i, v := range voivodeships {
		g.Go(func() error {
			data, err := s.AggregateForVoivodeship(ctx, v)
			if err != nil {
				return fmt.Errorf("aggregating %s: %w", v, err)
			}
			results[i] = data
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Service) AggregateForVoivodeship(ctx context.Context, voivodeship api.Voivodeship) (api.AggregatedData, error) {
	if err := ctx.Err(); err != nil {
		return api.AggregatedData{}, fmt.Errorf("context cancelled before aggregation: %w", err)
	}

	c := s.readCache()
	if c.err != nil {
		return api.AggregatedData{}, fmt.Errorf("service initialization failed: %w", c.err)
	}

	g, ctx := errgroup.WithContext(ctx)

	var openMeteoAverages map[api.ParamType]float32
	var openAqAverages map[api.ParamType]float32

	g.Go(func() error {
		averages, err := s.calculateOpenMeteoAverages(ctx, c.openMeteoParameters, c.openMeteoMap[voivodeship])
		if err != nil {
			return fmt.Errorf("failed to calculate averages for open meteo parameters: %w", err)
		}
		openMeteoAverages = averages
		return nil
	})

	g.Go(func() error {
		averages, err := s.calculateOpenAqAverages(ctx, c.openaqParameters, c.openaqMap[voivodeship])
		if err != nil {
			return fmt.Errorf("failed to calculate averages for open aq parameters: %w", err)
		}
		openAqAverages = averages
		return nil
	})

	if err := g.Wait(); err != nil {
		return api.AggregatedData{}, err
	}

	results := api.AggregatedData{Voivodeship: voivodeship}
	if err := results.AddParamInfoFromOpenMeteo(c.openMeteoParameters); err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to add param info: %w", err)
	}
	results.AddParamValues(mergeAverages(openMeteoAverages, openAqAverages))
	return results, nil
}

func (s *Service) calculateOpenMeteoAverages(ctx context.Context, parameters []openmeteo.Parameter, stations []openmeteo.Station) (map[api.ParamType]float32, error) {
	results := make([][]openmeteo.Measurement, len(stations))

	g, ctx := errgroup.WithContext(ctx)
	for i, station := range stations {
		g.Go(func() error {
			m, err := s.openmeteoClient.GetMeasurementForStation(ctx, station.Id)
			if err != nil {
				return fmt.Errorf("fetching open meteo measurements for station %d: %w", station.Id, err)
			}
			results[i] = m
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	measurements := make([]openmeteo.Measurement, 0)
	for _, m := range results {
		measurements = append(measurements, m...)
	}
	parameterMap := buildOpenMeteoParameterMap(parameters)
	return calculateAverage(groupByParamId(measurements, parameterMap)), nil
}

func (s *Service) calculateOpenAqAverages(ctx context.Context, parameters []openaq.Parameter, stations []openaq.Station) (map[api.ParamType]float32, error) {
	results := make([][]openaq.Measurement, len(stations))

	g, ctx := errgroup.WithContext(ctx)
	for i, station := range stations {
		g.Go(func() error {
			m, err := s.openaqClient.GetMeasurementForStation(ctx, station.Id)
			if err != nil {
				return fmt.Errorf("fetching open aq measurements for station %d: %w", station.Id, err)
			}
			results[i] = m
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	measurements := make([]openaq.Measurement, 0)
	for _, m := range results {
		measurements = append(measurements, m...)
	}
	parameterMap := buildOpenAqParameterMap(parameters)
	return calculateAverage(groupByParamId(measurements, parameterMap)), nil
}

func buildOpenMeteoParameterMap(parameters []openmeteo.Parameter) map[int]api.ParamType {
	paramIdAndType := make(map[int]api.ParamType)
	for _, param := range parameters {
		pt, err := api.MapOpenMeteoParamName(param.Name)
		if err != nil {
			slog.Debug("Unsupported openmeteo parameter", "name", param.Name, "id", param.Id)
			continue
		}
		paramIdAndType[param.Id] = pt
	}
	return paramIdAndType
}

func buildOpenAqParameterMap(parameters []openaq.Parameter) map[int]api.ParamType {
	paramIdAndType := make(map[int]api.ParamType)
	for _, param := range parameters {
		pt, err := api.MapOpenAqParamName(param.Name)
		if err != nil {
			slog.Debug("Unsupported openaq parameter", "name", param.Name, "id", param.Id)
			continue
		}
		paramIdAndType[param.Id] = pt
	}
	return paramIdAndType
}

type measurable interface {
	GetParameterId() int
	GetValue() float32
}

func groupByParamId[T measurable](measurements []T, paramMap map[int]api.ParamType) map[api.ParamType][]T {
	grouped := make(map[api.ParamType][]T)
	for _, m := range measurements {
		paramType, exists := paramMap[m.GetParameterId()]
		if !exists {
			slog.Debug("Parameter ID not found in map", "parameterId", m.GetParameterId())
			continue
		}
		grouped[paramType] = append(grouped[paramType], m)
	}
	return grouped
}

func calculateAverage[T measurable](grouped map[api.ParamType][]T) map[api.ParamType]float32 {
	averages := make(map[api.ParamType]float32, len(grouped))
	for paramType, mList := range grouped {
		if len(mList) == 0 {
			continue
		}
		var sum float32 = 0.0
		for _, m := range mList {
			sum += m.GetValue()
		}
		averages[paramType] = sum / float32(len(mList))
	}
	return averages
}

func mergeAverages(openMeteoMap, openAqMap map[api.ParamType]float32) map[api.ParamType]float32 {
	result := make(map[api.ParamType]float32)
	for paramType, value1 := range openMeteoMap {
		if value2, exists := openAqMap[paramType]; exists {
			result[paramType] = (value1 + value2) / 2.0
		} else {
			result[paramType] = value1
		}
	}
	for paramType, value2 := range openAqMap {
		if _, exists := openMeteoMap[paramType]; !exists {
			result[paramType] = value2
		}
	}
	return result
}
