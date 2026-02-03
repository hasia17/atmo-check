package aggregator

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"fmt"
	"log/slog"
	"time"
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
		slog.Error("Error getting open meteo stations")
		return nil, err
	}
	return groupStationsByVoivodeship(stations)
}

func (s *Service) groupOpenaqStations() (Map[openaq.Station], error) {
	stations, err := s.openaqClient.GetStations()
	if err != nil {
		slog.Error("Error getting openaq stations")
		return nil, err
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
				slog.Info("Station assigned to Voivodeship",
					"station", s.StationName(),
					"Voivodeship", v)
				vm[v] = append(vm[v], s)
			}
		}
	}
	slog.Info("Assigned Voivodeship stations", "stations", vm)
	return vm, nil
}

func stationInVoivodeship[T locatable](s T, b geographicalBounds) bool {
	return s.Latitude() >= b.MinLatitude &&
		s.Latitude() <= b.MaxLatitude &&
		s.Longitude() >= b.MinLongitude &&
		s.Longitude() <= b.MaxLongitude
}

func (s *Service) AggregateOpenMeteo(voivodeship api.Voivodeship) (api.AggregatedData, error) {
	measurements := make([]openmeteo.Measurement, 0)
	mappedVoivodeship, err := mapVoivodeship(voivodeship)
	if err != nil {
		slog.Error("Error mapping open meteo Voivodeship: ", err)
		return api.AggregatedData{}, err
	}
	for _, station := range s.openMeteoMap[mappedVoivodeship] {
		m, err := s.openmeteoClient.GetMeasurementForStation(station.Id)
		if err != nil {
			slog.Error("Error getting open meteo measurements: ", err)
			return api.AggregatedData{}, err
		}
		measurements = append(measurements, m...)
	}
	ad := api.AggregatedData{
		Voivodeship: voivodeship,
		Parameters:  calculateAverage(groupByParamId(slice(measurements))),
	}
	result, err := s.addOpenAqParamInfo(ad)
	if err != nil {
		slog.Error("Error adding parameters info: ", err)
		return api.AggregatedData{}, err
	}
	return result, nil
}

func (s *Service) AggregateOpenaq(voivodeship api.Voivodeship) (api.AggregatedData, error) {
	measurements := make([]openaq.Measurement, 0)
	mappedVoivodeship, err := mapVoivodeship(voivodeship)
	if err != nil {
		slog.Error("Error mapping open aq Voivodeship: ", err)
		return api.AggregatedData{}, err
	}
	for _, station := range s.openaqMap[mappedVoivodeship] {
		m, err := s.openaqClient.GetMeasurementForStation(station.Id)
		if err != nil {
			slog.Error("Error getting open aq measurements: ", err)
			return api.AggregatedData{}, err
		}
		measurements = append(measurements, m...)
	}
	ad := api.AggregatedData{
		Voivodeship: voivodeship,
		Parameters:  calculateAverage(groupByParamId(slice(measurements))),
	}
	result, err := s.addOpenAqParamInfo(ad)
	if err != nil {
		slog.Error("Error adding parameters info: ", err)
		return api.AggregatedData{}, err
	}
	return result, nil
}

type measurable interface {
	GetParameterId() int
	GetValue() float32
}

func groupByParamId(measurements []measurable) map[int][]measurable {
	grouped := make(map[int][]measurable)

	for _, m := range measurements {
		grouped[m.GetParameterId()] = append(grouped[m.GetParameterId()], m)
	}
	return grouped
}

func calculateAverage(grouped map[int][]measurable) []api.Parameter {
	params := make([]api.Parameter, 0)
	for p, mList := range grouped {
		var sum float32 = 0.0
		for _, m := range mList {
			sum += m.GetValue()
		}
		param := api.Parameter{
			Id:    p,
			Value: sum / float32(len(mList)),
		}
		params = append(params, param)
	}
	return params
}

func slice[T measurable](measurements []T) []measurable {
	result := make([]measurable, len(measurements))
	for i, m := range measurements {
		result[i] = m
	}
	return result
}

type parameterizable interface {
	GetName() string
	GetUnit() string
	GetDescription() string
	GetId() int
}

func (s *Service) addOpenMeteoParamInfo(data api.AggregatedData) (api.AggregatedData, error) {
	parameters, err := s.openmeteoClient.GetParameters()
	if err != nil {
		return data, fmt.Errorf("failed to fetch parameters from open meteo: %w", err)
	}
	addInfo(&data, paramSlice(parameters))
	return data, nil
}

func addInfo(data *api.AggregatedData, parameters []parameterizable) {
	pMap := make(map[int]parameterizable, len(parameters))
	for _, param := range parameters {
		pMap[param.GetId()] = param
	}
	for i := range data.Parameters {
		if param, exists := pMap[data.Parameters[i].Id]; exists {
			data.Parameters[i].Name = param.GetName()
			data.Parameters[i].Unit = param.GetUnit()
			data.Parameters[i].Description = param.GetDescription()
		}
	}
	data.Timestamp = time.Now().UTC().Format(time.RFC3339)
}

func (s *Service) addOpenAqParamInfo(data api.AggregatedData) (api.AggregatedData, error) {
	parameters, err := s.openaqClient.GetParameters()
	if err != nil {
		return data, fmt.Errorf("failed to fetch parameters from aq meteo: %w", err)
	}
	addInfo(&data, paramSlice(parameters))
	return data, nil
}

func paramSlice[T parameterizable](parameters []T) []parameterizable {
	result := make([]parameterizable, len(parameters))
	for i, m := range parameters {
		result[i] = m
	}
	return result
}
