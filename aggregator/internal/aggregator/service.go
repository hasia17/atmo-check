package aggregator

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"fmt"
	"log/slog"
)

type Voivodeship string

const (
	Dolnoslaskie       Voivodeship = "dolnoslaskie"
	KujawskoPomorskie  Voivodeship = "kujawsko-pomorskie"
	Lubelskie          Voivodeship = "lubelskie"
	Lubuskie           Voivodeship = "lubuskie"
	Lodzkie            Voivodeship = "lodzkie"
	Malopolskie        Voivodeship = "malopolskie"
	Mazowieckie        Voivodeship = "mazowieckie"
	Opolskie           Voivodeship = "opolskie"
	Podkarpackie       Voivodeship = "podkarpackie"
	Podlaskie          Voivodeship = "podlaskie"
	Pomorskie          Voivodeship = "pomorskie"
	Slaskie            Voivodeship = "slaskie"
	Swietokrzyskie     Voivodeship = "swietokrzyskie"
	WarminskoMazurskie Voivodeship = "warminsko-mazurskie"
	Wielkopolskie      Voivodeship = "wielkopolskie"
	Zachodniopomorskie Voivodeship = "zachodniopomorskie"
)

type Service struct {
	openmeteoClient *openmeteo.Client
	openaqClient    *openaq.Client
	openMeteoMap    Map[openmeteo.Station]
	openaqMap       Map[openaq.Station]
}

func NewService(openmeteoClient *openmeteo.Client, openaqClient *openaq.Client) (*Service, error) {
	s := &Service{
		openmeteoClient: openmeteoClient,
		openaqClient:    openaqClient,
	}
	err := s.initStationsMapping()
	if err != nil {
		return nil, fmt.Errorf("failed to init stations maps: %w", err)
	}

	return s, nil
}

func (s *Service) initStationsMapping() error {
	openMeteoMap, err := s.groupOpenMeteoStations()
	if err != nil {

		return fmt.Errorf("failed to group open meteo stations: %w", err)
	}
	s.openMeteoMap = openMeteoMap

	openaqMap, err := s.groupOpenaqStations()
	if err != nil {
		return fmt.Errorf("failed to group open aq stations: %w", err)
	}
	s.openaqMap = openaqMap
	return nil
}

type Map[T any] map[Voivodeship][]T

func (s *Service) groupOpenMeteoStations() (Map[openmeteo.Station], error) {
	stations, err := s.openmeteoClient.GetStations()
	if err != nil {
		slog.Error("Failed to get openmeteo stations", "error", err)
		return nil, fmt.Errorf("fetching openmeteo stations: %w", err)
	}
	return groupStationsByVoivodeship(stations)
}

func (s *Service) groupOpenaqStations() (Map[openaq.Station], error) {
	stations, err := s.openaqClient.GetStations()
	if err != nil {
		slog.Error("Failed to get openaq stations", "error", err)
		return nil, fmt.Errorf("fetching openaq stations: %w", err)
	}
	return groupStationsByVoivodeship(stations)
}

type geographicalBounds struct {
	MaxLatitude  float64
	MinLatitude  float64
	MaxLongitude float64
	MinLongitude float64
}

func voivodeshipBounds() map[Voivodeship]geographicalBounds {

	return map[Voivodeship]geographicalBounds{
		Zachodniopomorskie: {MinLatitude: 52.6, MaxLatitude: 54.9, MinLongitude: 14.0, MaxLongitude: 16.5},
		Pomorskie:          {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 16.5, MaxLongitude: 19.0},
		WarminskoMazurskie: {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 19.0, MaxLongitude: 22.0},
		Podlaskie:          {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 22.0, MaxLongitude: 24.2},
		Lubuskie:           {MinLatitude: 51.0, MaxLatitude: 52.6, MinLongitude: 14.0, MaxLongitude: 16.5},
		Wielkopolskie:      {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 16.5, MaxLongitude: 18.5},
		KujawskoPomorskie:  {MinLatitude: 52.0, MaxLatitude: 53.3, MinLongitude: 18.5, MaxLongitude: 20.0},
		Mazowieckie:        {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 20.0, MaxLongitude: 22.5},
		Lubelskie:          {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 22.5, MaxLongitude: 24.2},
		Dolnoslaskie:       {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 14.0, MaxLongitude: 16.8},
		Opolskie:           {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 16.8, MaxLongitude: 18.2},
		Slaskie:            {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 18.2, MaxLongitude: 19.5},
		Malopolskie:        {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 19.5, MaxLongitude: 21.0},
		Podkarpackie:       {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 21.0, MaxLongitude: 24.2},
		Lodzkie:            {MinLatitude: 50.8, MaxLatitude: 52.0, MinLongitude: 18.9, MaxLongitude: 20.5},
		Swietokrzyskie:     {MinLatitude: 50.1, MaxLatitude: 51.3, MinLongitude: 20.5, MaxLongitude: 21.8},
	}
}

type locatable interface {
	Latitude() float64
	Longitude() float64
	StationName() string
}

func groupStationsByVoivodeship[T locatable](stations []T) (Map[T], error) {

	vm := make(map[Voivodeship][]T)

	for v, b := range voivodeshipBounds() {
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

func (s *Service) AggregateData(voivodeship api.Voivodeship) (api.AggregatedData, error) {
	results := api.AggregatedData{
		Voivodeship: voivodeship,
	}

	openMeteoParameters, err := s.openmeteoClient.GetParameters()
	if err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to fetch parameters from open meteo: %w", err)
	}
	openMeteoAverages, err := s.calculateOpenMeteoAverages(openMeteoParameters, voivodeship)
	if err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to calculate averages for open meteo parameters: %w", err)
	}
	openAqAverages, err := s.calculateOpenAqAverages(voivodeship)
	if err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to calculate averages for open aq parameters: %w", err)
	}
	aggregatedParameters := mergeAverages(openMeteoAverages, openAqAverages)
	// take additional param info only from open-meteo
	if err := results.AddParamInfoFromOpenMeteo(openMeteoParameters); err != nil {
		return api.AggregatedData{}, fmt.Errorf("failed to add param info: %w", err)
	}
	results.AddParamValues(aggregatedParameters)
	return results, nil
}

func (s *Service) calculateOpenMeteoAverages(parameters []openmeteo.Parameter, voivodeship api.Voivodeship) (map[api.ParamType]float32, error) {
	measurements := make([]openmeteo.Measurement, 0)
	mappedVoivodeship, err := mapVoivodeship(voivodeship)
	if err != nil {
		return nil, fmt.Errorf("mapping voivodeship: %w", err)
	}
	for _, station := range s.openMeteoMap[mappedVoivodeship] {
		m, err := s.openmeteoClient.GetMeasurementForStation(station.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching measurements for station %d: %w", station.Id, err)
		}
		measurements = append(measurements, m...)
	}
	parameterMap := buildOpenMeteoParameterMap(parameters)
	return calculateAverage(groupByParamId(measurements, parameterMap)), nil
}

func (s *Service) calculateOpenAqAverages(voivodeship api.Voivodeship) (map[api.ParamType]float32, error) {
	parameters, err := s.openaqClient.GetParameters()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parameters from open aq: %w", err)
	}
	measurements := make([]openaq.Measurement, 0)
	mappedVoivodeship, err := mapVoivodeship(voivodeship)
	if err != nil {
		return nil, fmt.Errorf("mapping voivodeship: %w", err)
	}
	for _, station := range s.openaqMap[mappedVoivodeship] {
		m, err := s.openaqClient.GetMeasurementForStation(station.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching measurements for station %d: %w", station.Id, err)
		}
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
