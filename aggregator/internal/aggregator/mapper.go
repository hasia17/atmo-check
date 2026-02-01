package aggregator

import (
	"aggregator/internal/api"
	"fmt"
)

func mapVoivodeship(v api.Voivodeship) (Voivodeship, error) {
	switch v {
	case api.Dolnoslaskie:
		return Dolnoslaskie, nil
	case api.KujawskoPomorskie:
		return KujawskoPomorskie, nil
	case api.Lubelskie:
		return Lubelskie, nil
	case api.Lubuskie:
		return Lubuskie, nil
	case api.Lodzkie:
		return Lodzkie, nil
	case api.Malopolskie:
		return Malopolskie, nil
	case api.Mazowieckie:
		return Mazowieckie, nil
	case api.Opolskie:
		return Opolskie, nil
	case api.Podkarpackie:
		return Podkarpackie, nil
	case api.Podlaskie:
		return Podlaskie, nil
	case api.Pomorskie:
		return Pomorskie, nil
	case api.Slaskie:
		return Slaskie, nil
	case api.Swietokrzyskie:
		return Swietokrzyskie, nil
	case api.WarminskoMazurskie:
		return WarminskoMazurskie, nil
	case api.Wielkopolskie:
		return Wielkopolskie, nil
	case api.Zachodniopomorskie:
		return Zachodniopomorskie, nil
	default:
		return "", fmt.Errorf("unknown voivodeship: %s", v)
	}
}
