package api

type AggregatedData struct {
	Voivodeship Voivodeship `json:"voivodeship"`
	Parameters  []Parameter `json:"parameters"`
	Timestamp   string      `json:"timestamp"`
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
