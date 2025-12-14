package apiclient

type OpenAqLocation struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Locality string `json:"locality"`
	Timezone string `json:"timezone"`
	Country  struct {
		Id   int32  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`
	Sensors []struct {
		Id        int32  `json:"id"`
		Name      string `json:"name"`
		Parameter struct {
			Id          int32  `json:"id"`
			Name        string `json:"name"`
			Units       string `json:"units"`
			DisplayName string `json:"displayName"`
		} `json:"parameter"`
	} `json:"sensors"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
}

type openAqLocationResponse struct {
	Results []OpenAqLocation `json:"results"`
}

type OpenAqMeasurement struct {
	Date struct {
		Utc   string `json:"utc"`
		Local string `json:"local"`
	} `json:"datetime"`
	Value       float64 `json:"value"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	SensorId   int32 `json:"sensorsId"`
	LocationId int32 `json:"locationsId"`
}

type openAqMeasurementResponse struct {
	Results []OpenAqMeasurement `json:"results"`
}

type OpenAqParameter struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Units       string `json:"units"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type openAqParameterResponse struct {
	Results []OpenAqParameter `json:"results"`
}
