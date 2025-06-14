package internal

import "time"

type Station struct {
	ID         int32       `json:"id"                   bson:"_id"`
	Name       string      `json:"name"                 bson:"name"`
	Locality   string      `json:"locality,omitempty"   bson:"locality,omitempty"`
	Timezone   string      `json:"timezone,omitempty"   bson:"timezone,omitempty"`
	Country    Country     `json:"country"              bson:"country,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty" bson:"parameters,omitempty"`
}

type Country struct {
	ID   int32  `json:"id"   bson:"id"`
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

type Parameter struct {
	ID          int32  `json:"id"          bson:"id"`
	Name        string `json:"name"        bson:"name"`
	Units       string `json:"units"       bson:"units"`
	DisplayName string `json:"displayName" bson:"displayName"`
}

type Measurement struct {
	DateTime    MeasurementDateTime `json:"datetime"    bson:"datetime"`
	Timestamp   time.Time           `json:"-"           bson:"timestamp"`
	Value       float64             `json:"value"       bson:"value"`
	Coordinates Coordinates         `json:"coordinates" bson:"coordinates"`
	SensorID    int32               `json:"sensorId"    bson:"sensorId"`
	StationID   int32               `json:"stationId"   bson:"stationId"`
}

type MeasurementDateTime struct {
	UTC   string `json:"utc"   bson:"utc"`
	Local string `json:"local" bson:"local"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"  bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}
