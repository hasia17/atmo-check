package openaq

import "time"

type Measurement struct {
	ParameterId int       `json:"parameterId"`
	StationId   int       `json:"stationId"`
	Timestamp   time.Time `json:"timestamp"`
	Value       float32   `json:"value"`
}

type Parameter struct {
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Units       string `json:"units"`
}

type Station struct {
	Id           int     `json:"id"`
	Lat          float64 `json:"latitude"`
	Locality     string  `json:"locality"`
	Lon          float64 `json:"longitude"`
	Name         string  `json:"name"`
	ParameterIds []int   `json:"parameterIds"`
	Timezone     string  `json:"timezone"`
}

func (s Station) Longitude() float64 {
	return s.Lon
}

func (s Station) Latitude() float64 {
	return s.Lat
}

func (s Station) StationName() string {
	return s.Name
}

func (m Measurement) GetParameterId() int { return m.ParameterId }

func (m Measurement) GetValue() float32 { return m.Value }
