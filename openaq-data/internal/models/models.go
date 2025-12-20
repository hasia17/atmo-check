package models

// Package models provides data types that represent OpenAQ API responses.
type Location struct {
	Id       int32  `json:"id"          bson:"id"`
	Name     string `json:"name"        bson:"name"`
	Locality string `json:"locality"    bson:"locality"`
	Timezone string `json:"timezone"    bson:"timezone"`
	Country  struct {
		Id   int32  `json:"id" bson:"id"`
		Code string `json:"code" bson:"code"`
		Name string `json:"name" bson:"name"`
	} `json:"country"     bson:"country"`
	Sensors     []Sensor `json:"sensors"     bson:"sensors"`
	Coordinates struct {
		Latitude  float64 `json:"latitude" bson:"latitude"`
		Longitude float64 `json:"longitude" bson:"longitude"`
	} `json:"coordinates" bson:"coordinates"`
}

type Sensor struct {
	Id        int32     `json:"id"        bson:"id"`
	Name      string    `json:"name"      bson:"name"`
	Parameter Parameter `json:"parameter" bson:"parameter"`
}

type Parameter struct {
	Id          int32  `json:"id"          bson:"id"`
	Name        string `json:"name"        bson:"name"`
	Units       string `json:"units"       bson:"units"`
	DisplayName string `json:"displayName" bson:"displayName"`
	Description string `json:"description" bson:"description"`
}

type Measurement struct {
	Date struct {
		Utc   string `json:"utc" bson:"utc"`
		Local string `json:"local" bson:"local"`
	} `json:"datetime"    bson:"datetime"`
	Value       float64 `json:"value"       bson:"value"`
	Coordinates struct {
		Latitude  float64 `json:"latitude" bson:"latitude"`
		Longitude float64 `json:"longitude" bson:"longitude"`
	} `json:"coordinates" bson:"coordinates"`
	SensorId   int32 `json:"sensorsId"   bson:"sensorsId"`
	LocationId int32 `json:"locationsId" bson:"locationsId"`
}
