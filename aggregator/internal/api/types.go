package api

import (
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
)

type AggregatedData struct {
	Voivodeship Voivodeship `json:"voivodeship"`
	Parameters  []Parameter `json:"parameters"`
	Timestamp   string      `json:"timestamp"`
}

func (ad *AggregatedData) AddOpenAqParamInfo(parameters []openaq.Parameter) {
	pMap := make(map[int]openaq.Parameter, len(parameters))
	for _, param := range parameters {
		pMap[param.Id] = param
	}
	for i := range ad.Parameters {
		if param, exists := pMap[ad.Parameters[i].Id]; exists {
			ad.Parameters[i].Name = param.Name
			ad.Parameters[i].Unit = param.Units
			ad.Parameters[i].Description = param.Description
		}
	}
}

func (ad *AggregatedData) AddOpenMeteoParamInfo(parameters []openmeteo.Parameter) {
	pMap := make(map[int]openmeteo.Parameter, len(parameters))
	for _, param := range parameters {
		pMap[param.Id] = param
	}
	for i := range ad.Parameters {
		if param, exists := pMap[ad.Parameters[i].Id]; exists {
			ad.Parameters[i].Name = param.Name
			ad.Parameters[i].Unit = param.Unit
			ad.Parameters[i].Description = param.Description
		}
	}
}

type Parameter struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Value       float32 `json:"value"`
}

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
