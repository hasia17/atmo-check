package aggregator

import (
	"aggregator/internal/api"
	"fmt"
)

func mapVoivodeship(v Voivodeship) (api.Voivodeship, error) {
	switch v {
	case Dolnoslaskie:
		return api.Dolnoslaskie, nil
	case KujawskoPomorskie:
		return api.KujawskoPomorskie, nil
	case Lubelskie:
		return api.Lubelskie, nil
	case Lubuskie:
		return api.Lubuskie, nil
	case Lodzkie:
		return api.Lodzkie, nil
	case Malopolskie:
		return api.Malopolskie, nil
	case Mazowieckie:
		return api.Mazowieckie, nil
	case Opolskie:
		return api.Opolskie, nil
	case Podkarpackie:
		return api.Podkarpackie, nil
	case Podlaskie:
		return api.Podlaskie, nil
	case Pomorskie:
		return api.Pomorskie, nil
	case Slaskie:
		return api.Slaskie, nil
	case Swietokrzyskie:
		return api.Swietokrzyskie, nil
	case WarminskoMazurskie:
		return api.WarminskoMazurskie, nil
	case Wielkopolskie:
		return api.Wielkopolskie, nil
	case Zachodniopomorskie:
		return api.Zachodniopomorskie, nil
	default:
		return "", fmt.Errorf("unknown voivodeship: %s", v)
	}
}
