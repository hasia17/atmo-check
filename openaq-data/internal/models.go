package internal

import "time"

// Location represents a location with air quality sensors.
type Location struct {
	ID       int32    `json:"id"                 bson:"_id"`
	Name     string   `json:"name"               bson:"name"`
	Locality string   `json:"locality,omitempty" bson:"locality,omitempty"`
	Timezone string   `json:"timezone,omitempty" bson:"timezone,omitempty"`
	Country  Country  `json:"country"            bson:"country,omitempty"`
	Sensors  []Sensor `json:"sensors,omitempty"  bson:"sensors,omitempty"`
}

// Country represents a country.
type Country struct {
	ID   int32  `json:"id"   bson:"id"`
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

// Sensor represents a sensor at a location.
type Sensor struct {
	ID        int32     `json:"id"        bson:"_id"`
	Name      string    `json:"name"      bson:"name"`
	Parameter Parameter `json:"parameter" bson:"parameter"`
}

// Parameter represents a measured parameter.
type Parameter struct {
	ID          int32  `json:"id"          bson:"id"`
	Name        string `json:"name"        bson:"name"`
	Units       string `json:"units"       bson:"units"`
	DisplayName string `json:"displayName" bson:"displayName"`
}

// Measurement represents an air quality measurement.
type Measurement struct {
	DateTime    MeasurementDateTime `json:"datetime"    bson:"datetime"`
	Timestamp   time.Time           `json:"-"           bson:"timestamp"`
	Value       float64             `json:"value"       bson:"value"`
	Coordinates Coordinates         `json:"coordinates" bson:"coordinates"`
	SensorID    int32               `json:"sensorId"    bson:"sensorId"`
	LocationID  int32               `json:"locationId"  bson:"locationId"`
}

// MeasurementDateTime represents the date and time of a measurement.
type MeasurementDateTime struct {
	UTC   string `json:"utc"   bson:"utc"`
	Local string `json:"local" bson:"local"`
}

// Coordinates represents geographical coordinates.
type Coordinates struct {
	Latitude  float64 `json:"latitude"  bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}
