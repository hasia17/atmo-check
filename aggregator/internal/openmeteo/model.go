package openmeteo

type Station struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	GeoLat float64 `json:"geoLat"`
	GeoLon float64 `json:"geoLon"`
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

func (s Station) GetLatitude() float64 {
	return s.GeoLat
}

func (s Station) GetLongitude() float64 {
	return s.GeoLon
}

func (s Station) GetName() string {
	return s.Name
}
