package service

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"log"
)

func GroupOpenMeteoStations() (map[api.Voivodeship][]openmeteo.Station, error) {

	stations, err := openmeteo.GetStations()
	if err != nil {
		log.Print("Error getting open meteo stations")
		return nil, err
	}

	return GroupStationsByVoivodeship(stations)
}

func GroupOpenAqStations() (map[api.Voivodeship][]openaq.Station, error) {

	stations, err := openaq.GetStations()
	if err != nil {
		log.Print("Error getting openaq stations")
		return nil, err
	}

	return GroupStationsByVoivodeship(stations)
}

func getVoivodeshipBounds() map[api.Voivodeship]api.GeographicalBounds {

	return map[api.Voivodeship]api.GeographicalBounds{
		api.Zachodniopomorskie: {MinLatitude: 52.6, MaxLatitude: 54.9, MinLongitude: 14.0, MaxLongitude: 16.5},
		api.Pomorskie:          {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 16.5, MaxLongitude: 19.0},
		api.WarminskoMazurskie: {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 19.0, MaxLongitude: 22.0},
		api.Podlaskie:          {MinLatitude: 53.3, MaxLatitude: 54.9, MinLongitude: 22.0, MaxLongitude: 24.2},
		api.Lubuskie:           {MinLatitude: 51.0, MaxLatitude: 52.6, MinLongitude: 14.0, MaxLongitude: 16.5},
		api.Wielkopolskie:      {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 16.5, MaxLongitude: 18.5},
		api.KujawskoPomorskie:  {MinLatitude: 52.0, MaxLatitude: 53.3, MinLongitude: 18.5, MaxLongitude: 20.0},
		api.Mazowieckie:        {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 20.0, MaxLongitude: 22.5},
		api.Lubelskie:          {MinLatitude: 51.0, MaxLatitude: 53.3, MinLongitude: 22.5, MaxLongitude: 24.2},
		api.Dolnoslaskie:       {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 14.0, MaxLongitude: 16.8},
		api.Opolskie:           {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 16.8, MaxLongitude: 18.2},
		api.Slaskie:            {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 18.2, MaxLongitude: 19.5},
		api.Malopolskie:        {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 19.5, MaxLongitude: 21.0},
		api.Podkarpackie:       {MinLatitude: 49.0, MaxLatitude: 51.0, MinLongitude: 21.0, MaxLongitude: 24.2},
		api.Lodzkie:            {MinLatitude: 50.8, MaxLatitude: 52.0, MinLongitude: 18.9, MaxLongitude: 20.5},
		api.Swietokrzyskie:     {MinLatitude: 50.1, MaxLatitude: 51.3, MinLongitude: 20.5, MaxLongitude: 21.8},
	}
}

func GroupStationsByVoivodeship[T api.StationWithCoordinates](stations []T) (map[api.Voivodeship][]T, error) {

	stationsByVoivodeship := make(map[api.Voivodeship][]T)

	voivodeshipBounds := getVoivodeshipBounds()

	for voivodeship, bounds := range voivodeshipBounds {

		for _, station := range stations {

			if station.GetLatitude() >= bounds.MinLatitude &&
				station.GetLatitude() <= bounds.MaxLatitude &&
				station.GetLongitude() >= bounds.MinLongitude &&
				station.GetLongitude() <= bounds.MaxLongitude {

				log.Printf("Stations %v assigned to voivodeship: %s", station.GetName(), voivodeship)
				stationsByVoivodeship[voivodeship] = append(stationsByVoivodeship[voivodeship], station)
			}
		}
	}
	log.Print("Assigned voivodeship stations: ", stationsByVoivodeship)
	return stationsByVoivodeship, nil
}
