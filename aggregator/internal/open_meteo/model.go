package open_meteo

type Station struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	GeoLat  float64 `json:"geoLat"`
	GeoLong float64 `json:"geoLong"`
}

type Parameter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type Measurement struct {
	Id          string  `json:"id"`
	StationId   int     `json:"stationId"`
	ParameterId int     `json:"parameterId"`
	Value       float32 `json:"value"`
	Timestamp   string  `json:"timestamp"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
