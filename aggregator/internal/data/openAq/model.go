package openAq

import "time"

type Measurement struct {
	ParameterId int32     `json:"parameterId"`
	StationId   int32     `json:"stationId"`
	Timestamp   time.Time `json:"timestamp"`
	Value       float64   `json:"value"`
}

type Parameter struct {
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName"`
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Units       string `json:"units"`
}

type Station struct {
	Id           int32   `json:"id"`
	Latitude     float64 `json:"latitude"`
	Locality     string  `json:"locality"`
	Longitude    float64 `json:"longitude"`
	Name         string  `json:"name"`
	ParameterIds []int32 `json:"parameterIds"`
	Timezone     string  `json:"timezone"`
}
