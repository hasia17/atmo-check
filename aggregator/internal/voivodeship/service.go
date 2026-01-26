package voivodeship

import (
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"log/slog"
)

type Service struct {
	openmeteoClient *openmeteo.Client
	openaqClient    *openaq.Client
}

func NewService(openmeteoClient *openmeteo.Client, openaqClient *openaq.Client) *Service {
	return &Service{
		openmeteoClient: openmeteoClient,
		openaqClient:    openaqClient,
	}
}

type Map[T any] map[voivodeship][]T

func (s *Service) GroupOpenMeteoStations() (Map[openmeteo.Station], error) {

	stations, err := s.openmeteoClient.GetStations()
	if err != nil {
		slog.Info("Error getting open meteo stations")
		return nil, err
	}

	return groupStationsByVoivodeship(stations)
}

func (s *Service) GroupOpenAqStations() (Map[openaq.Station], error) {

	stations, err := s.openaqClient.GetStations()
	if err != nil {
		slog.Info("Error getting openaq stations")
		return nil, err
	}

	return groupStationsByVoivodeship(stations)
}

type voivodeship string

const (
	Dolnoslaskie       voivodeship = "dolnoslaskie"
	KujawskoPomorskie  voivodeship = "kujawsko-pomorskie"
	Lubelskie          voivodeship = "lubelskie"
	Lubuskie           voivodeship = "lubuskie"
	Lodzkie            voivodeship = "lodzkie"
	Malopolskie        voivodeship = "malopolskie"
	Mazowieckie        voivodeship = "mazowieckie"
	Opolskie           voivodeship = "opolskie"
	Podkarpackie       voivodeship = "podkarpackie"
	Podlaskie          voivodeship = "podlaskie"
	Pomorskie          voivodeship = "pomorskie"
	Slaskie            voivodeship = "slaskie"
	Swietokrzyskie     voivodeship = "swietokrzyskie"
	WarminskoMazurskie voivodeship = "warminsko-mazurskie"
	Wielkopolskie      voivodeship = "wielkopolskie"
	Zachodniopomorskie voivodeship = "zachodniopomorskie"
)

type geographicalBounds struct {
	MaxLatitude  float64
	MinLatitude  float64
	MaxLongitude float64
	MinLongitude float64
}

type locatable interface {
	Latitude() float64
	Longitude() float64
	StationName() string
}

func voivodeshipBounds() map[voivodeship]geographicalBounds {

	return map[voivodeship]geographicalBounds{
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

func groupStationsByVoivodeship[T locatable](stations []T) (Map[T], error) {

	vm := make(map[voivodeship][]T)

	for v, b := range voivodeshipBounds() {
		for _, s := range stations {
			if stationInVoivodeship(s, b) {
				slog.Info("Station assigned to voivodeship",
					"station", s.StationName(),
					"voivodeship", v)
				vm[v] = append(vm[v], s)
			}
		}
	}
	slog.Info("Assigned voivodeship stations", "stations", vm)
	return vm, nil
}

func stationInVoivodeship[T locatable](s T, b geographicalBounds) bool {
	return s.Latitude() >= b.MinLatitude &&
		s.Latitude() <= b.MaxLatitude &&
		s.Longitude() >= b.MinLongitude &&
		s.Longitude() <= b.MaxLongitude
}
