package api

import (
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"fmt"
	"slices"
	"time"
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

func (ad *AggregatedData) AddParamInfo(parameters []openmeteo.Parameter, averages map[ParamType]float32) error {
	validParamTypes := []string{"PM10", "PM2_5", "CARBON_MONOXIDE", "CARBON_DIOXIDE", "NITROGEN_DIOXIDE", "SULPHUR_DIOXIDE", "OZONE", "METHANE"}
	params := make([]Parameter, 0)

	for i, p := range parameters {
		if slices.Contains(validParamTypes, p.Name) {
			param, err := MapParameter(p)
			if err != nil {
				return fmt.Errorf("unknown parameter: %s", p)
			}
			param.Id = i
			if value, exists := averages[param.Type]; exists {
				param.Value = value
			}
			params = append(params, param)
		}
	}
	ad.Parameters = params
	ad.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return nil
}

type Parameter struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
	Value       float32   `json:"value"`
	Type        ParamType `json:"type"`
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
	PM10  ParamType = "pm10"
	PM2_5 ParamType = "pm2_5"
	CO    ParamType = "co"
	CO2   ParamType = "co2"
	NO2   ParamType = "no2"
	SO2   ParamType = "so2"
	O3    ParamType = "o3"
	CH4   ParamType = "ch4"
)
