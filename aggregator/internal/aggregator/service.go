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

	"golang.org/x/sync/errgroup"
)

type Service struct {
	openmeteoClient   *openmeteo.Client
	openaqClient      *openaq.Client
	openMeteoMap      Map[openmeteo.Station]
	openaqMap         Map[openaq.Station]
	voivodeshipBounds map[api.Voivodeship]geographicalBounds
}

func NewService(ctx context.Context) (*Service, error) {
	bounds, err := loadVoivodeshipBounds("config/voivodeships.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load voivodeship bounds: %w", err)
	}
	s := &Service{
		openmeteoClient:   openmeteo.NewClient(),
		openaqClient:      openaq.NewClient(),
		voivodeshipBounds: bounds,
	}
	err = s.initStationsMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init stations maps: %w", err)
	}

	return s, nil
}

func (s *Service) initStationsMapping(ctx context.Context) error {
	openMeteoMap, err := s.groupOpenMeteoStations(ctx)
	if err != nil {

		return fmt.Errorf("failed to group open meteo stations: %w", err)
	}
	s.openMeteoMap = openMeteoMap

	openaqMap, err := s.groupOpenaqStations(ctx)
	if err != nil {
		return fmt.Errorf("failed to group open aq stations: %w", err)
	}
	s.openaqMap = openaqMap
	return nil
}

type Map[T any] map[api.Voivodeship][]T

func (s *Service) groupOpenMeteoStations(ctx context.Context) (Map[openmeteo.Station], error) {
	stations, err := s.openmeteoClient.GetStations(ctx)
	if err != nil {
		slog.Error("Failed to get openmeteo stations", "error", err)
		return nil, fmt.Errorf("fetching openmeteo stations: %w", err)
	}
	return groupStationsByVoivodeship(stations, s.voivodeshipBounds)
}

func (s *Service) groupOpenaqStations(ctx context.Context) (Map[openaq.Station], error) {
	stations, err := s.openaqClient.GetStations(ctx)
	if err != nil {
		slog.Error("Failed to get openaq stations", "error", err)
		return nil, fmt.Errorf("fetching openaq stations: %w", err)
	}
	return groupStationsByVoivodeship(stations, s.voivodeshipBounds)
}

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

func groupStationsByVoivodeship[T locatable](stations []T, bounds map[api.Voivodeship]geographicalBounds) (Map[T], error) {
	vm := make(map[api.Voivodeship][]T)

	for v, b := range bounds {
		for _, s := range stations {
			if stationInVoivodeship(s, b) {
				vm[v] = append(vm[v], s)
			}
		}
	}
	return vm, nil
}

func stationInVoivodeship[T locatable](s T, b geographicalBounds) bool {
	return s.Latitude() >= b.MinLatitude &&
		s.Latitude() <= b.MaxLatitude &&
		s.Longitude() >= b.MinLongitude &&
		s.Longitude() <= b.MaxLongitude
}

func (s *Service) AggregateData(ctx context.Context, voivodeship api.Voivodeship) (api.AggregatedData, error) {
	if err := ctx.Err(); err != nil {
		return api.AggregatedData{}, fmt.Errorf("context cancelled before aggregation: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	var openMeteoParameters []openmeteo.Parameter
	var openMeteoAverages map[api.ParamType]float32
	var openAqAverages map[api.ParamType]float32

	g.Go(func() error {
		params, err := s.openmeteoClient.GetParameters(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch parameters from open meteo: %w", err)
		}
		averages, err := s.calculateOpenMeteoAverages(ctx, params, voivodeship)
		if err != nil {
			return fmt.Errorf("failed to calculate averages for open meteo parameters: %w", err)
		}
		openMeteoParameters = params
		openMeteoAverages = averages
		return nil
	})

	g.Go(func() error {
		averages, err := s.calculateOpenAqAverages(ctx, voivodeship)
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
	// take additional param info only from open-meteo
	if err := results.AddParamInfoFromOpenMeteo(openMeteoParameters); err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to add param info: %w", err)
	}
	results.AddParamValues(mergeAverages(openMeteoAverages, openAqAverages))
	return results, nil
}

func (s *Service) calculateOpenMeteoAverages(ctx context.Context, parameters []openmeteo.Parameter, voivodeship api.Voivodeship) (map[api.ParamType]float32, error) {
	stations := s.openMeteoMap[voivodeship]
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

func (s *Service) calculateOpenAqAverages(ctx context.Context, voivodeship api.Voivodeship) (map[api.ParamType]float32, error) {
	parameters, err := s.openaqClient.GetParameters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parameters from open aq: %w", err)
	}
	stations := s.openaqMap[voivodeship]
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
