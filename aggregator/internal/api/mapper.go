package api

import (
	"aggregator/internal/openmeteo"
	"fmt"
	"strings"
)

func MapOpenMeteoParameter(parameter openmeteo.Parameter) (Parameter, error) {
	paramType, err := MapOpenMeteoParamName(parameter.Name)
	if err != nil {
		return Parameter{}, fmt.Errorf("unsupported paramName: %s", parameter.Name)
	}
	return Parameter{
		Description: parameter.Description,
		Unit:        parameter.Unit,
		Type:        paramType,
		Id:          validParamTypes[paramType],
	}, nil
}

func MapOpenMeteoParamName(paramName string) (ParamType, error) {
	switch paramName {
	case "PM10":
		return PM10, nil
	case "PM2_5":
		return PM2_5, nil
	case "CARBON_MONOXIDE":
		return CO, nil
	case "CARBON_DIOXIDE":
		return CO2, nil
	case "NITROGEN_DIOXIDE":
		return NO2, nil
	case "SULPHUR_DIOXIDE":
		return SO2, nil
	case "OZONE":
		return O3, nil
	case "METHANE":
		return CH4, nil
	default:
		return "", fmt.Errorf("unsupported paramName: %s", paramName)
	}
}

func MapOpenAqParamName(paramName string) (ParamType, error) {
	switch paramName {
	case "pm10":
		return PM10, nil
	case "pm25":
		return PM2_5, nil
	case "co":
		return CO, nil
	case "co2":
		return CO2, nil
	case "no2":
		return NO2, nil
	case "so2":
		return SO2, nil
	case "o3":
		return O3, nil
	case "ch4":
		return CH4, nil
	default:
		return "", fmt.Errorf("unsupported paramName: %s", paramName)
	}
}

func MapVoivodeship(s string) (Voivodeship, error) {
	v := Voivodeship(strings.ToLower(s))
	switch v {
	case Dolnoslaskie, KujawskoPomorskie, Lubelskie, Lubuskie, Lodzkie,
		Malopolskie, Mazowieckie, Opolskie, Podkarpackie, Podlaskie,
		Pomorskie, Slaskie, Swietokrzyskie, WarminskoMazurskie,
		Wielkopolskie, Zachodniopomorskie:
		return v, nil
	default:
		return "", fmt.Errorf("unknown voivodeship: %s", s)
	}
}
