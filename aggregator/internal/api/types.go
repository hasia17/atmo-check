package api

import (
	"aggregator/internal/openmeteo"
	"time"
)

var validParamTypes = map[ParamType]int{
	"PM10":  1,
	"PM2_5": 2,
	"CO":    3,
	"CO2":   4,
	"NO2":   5,
	"SO2":   6,
	"O3":    7,
	"CH4":   8,
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

type ParamType string

const (
	PM10  ParamType = "PM10"
	PM2_5 ParamType = "PM2_5"
	CO    ParamType = "CO"
	CO2   ParamType = "CO2"
	NO2   ParamType = "NO2"
	SO2   ParamType = "SO2"
	O3    ParamType = "O3"
	CH4   ParamType = "CH4"
)

type Parameter struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
	Value       float32   `json:"value"`
	Type        ParamType `json:"type"`
}

type AggregatedData struct {
	Voivodeship Voivodeship `json:"voivodeship"`
	Parameters  []Parameter `json:"parameters"`
	Timestamp   string      `json:"timestamp"`
}

func (ad *AggregatedData) AddParamInfoFromOpenMeteo(parameters []openmeteo.Parameter) error {
	params := make([]Parameter, len(validParamTypes))
	for _, p := range parameters {
		param, err := MapOpenMeteoParameter(p)
		if err != nil {
			continue
		}
		params[param.Id-1] = param
	}
	ad.Parameters = params
	return nil
}

func (ad *AggregatedData) AddParamValues(averages map[ParamType]float32) {
	for i := range ad.Parameters {
		p := &ad.Parameters[i]
		if value, exists := averages[p.Type]; exists {
			ad.Parameters[i].Value = value
		}
	}
	ad.Timestamp = time.Now().UTC().Format(time.RFC3339)
}
