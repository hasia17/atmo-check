package aggregator

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"testing"
)

func TestBuildOpenMeteoParameterMap(t *testing.T) {
	p1 := openmeteo.Parameter{
		Id:   2,
		Name: "PM10",
	}
	p2 := openmeteo.Parameter{
		Id:   4,
		Name: "CARBON_MONOXIDE",
	}
	p3 := openmeteo.Parameter{
		Id:   6,
		Name: "UNKNOWN_PARAM",
	}
	params := []openmeteo.Parameter{p1, p2, p3}
	result := buildOpenMeteoParameterMap(params)
	if result[2] != api.PM10 {
		t.Errorf("Expected PM10, got %s", result[2])
	}
	if result[4] != api.CO {
		t.Errorf("Expected CO, got %s", result[4])
	}
	if _, exists := result[6]; exists {
		t.Error("Expected unknown parameter to be skipped")
	}
}

func TestGroupStationsByVoivodeship(t *testing.T) {
	b1 := geographicalBounds{
		MaxLatitude:  10,
		MinLatitude:  5,
		MaxLongitude: 10,
		MinLongitude: 5,
	}
	b2 := geographicalBounds{
		MaxLatitude:  20,
		MinLatitude:  11,
		MaxLongitude: 20,
		MinLongitude: 11,
	}
	b3 := geographicalBounds{
		MaxLatitude:  30,
		MinLatitude:  21,
		MaxLongitude: 30,
		MinLongitude: 21,
	}

	s1 := openmeteo.Station{
		GeoLat: 8,
		GeoLon: 10,
	}
	s2 := openmeteo.Station{
		GeoLat: 18,
		GeoLon: 15,
	}
	s3 := openmeteo.Station{
		GeoLat: 13,
		GeoLon: 15,
	}
	bounds := map[api.Voivodeship]geographicalBounds{
		api.Malopolskie: b1,
		api.Mazowieckie: b2,
		api.Pomorskie:   b3,
	}
	stations := []openmeteo.Station{s1, s2, s3}
	result := groupStationsByVoivodeship(stations, bounds)
	if len(result[api.Malopolskie]) != 1 {
		t.Errorf("Expected 1 station for Malopolskie, got %d", len(result[api.Malopolskie]))
	}
	if len(result[api.Mazowieckie]) != 2 {
		t.Errorf("Expected 1 station for Mazowieckie, got %d", len(result[api.Mazowieckie]))
	}
	if len(result[api.Pomorskie]) != 0 {
		t.Errorf("Expected 1 station for Pomorskie, got %d", len(result[api.Pomorskie]))
	}
}

func TestStationInVoivodeship(t *testing.T) {
	bounds := geographicalBounds{
		MaxLatitude:  10,
		MinLatitude:  5,
		MaxLongitude: 20,
		MinLongitude: 5,
	}

	station := openmeteo.Station{
		GeoLat: 8,
		GeoLon: 10,
	}
	result := stationInVoivodeship(station, bounds)
	if result != true {
		t.Errorf("Expected true, got %t", result)
	}
	station.GeoLat = 20
	result = stationInVoivodeship(station, bounds)
	if result != false {
		t.Errorf("Expected false, got %t", result)
	}
	station.GeoLon = 2
	result = stationInVoivodeship(station, bounds)
	if result != false {
		t.Errorf("Expected false, got %t", result)
	}
}

func TestBuildOpenAqParameterMap(t *testing.T) {
	p1 := openaq.Parameter{
		Id:   2,
		Name: "pm10",
	}
	p2 := openaq.Parameter{
		Id:   4,
		Name: "co",
	}
	p3 := openaq.Parameter{
		Id:   6,
		Name: "UNKNOWN_PARAM",
	}
	params := []openaq.Parameter{p1, p2, p3}
	result := buildOpenAqParameterMap(params)
	if result[2] != api.PM10 {
		t.Errorf("Expected PM10, got %s", result[2])
	}
	if result[4] != api.CO {
		t.Errorf("Expected CO, got %s", result[4])
	}
	if _, exists := result[6]; exists {
		t.Error("Expected unknown parameter to be skipped")
	}
}

func TestGroupByParamId(t *testing.T) {
	paramMap := map[int]api.ParamType{
		1: api.PM10,
		2: api.SO2,
	}

	measurements := []openmeteo.Measurement{
		{ParameterId: 1, Value: 10},
		{ParameterId: 1, Value: 20},
		{ParameterId: 2, Value: 30},
		{ParameterId: 99, Value: 40},
	}

	result := groupByParamId(measurements, paramMap)

	if len(result[api.PM10]) != 2 {
		t.Errorf("Expected 2 measurements for PM10, got %d", len(result[api.PM10]))
	}
	if len(result[api.SO2]) != 1 {
		t.Errorf("Expected 1 measurement for SO2, got %d", len(result[api.SO2]))
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}
}

func TestMergeAverages(t *testing.T) {
	openMeteo := map[api.ParamType]float32{
		api.SO2: 50,
		api.CH4: 30,
	}
	openAq := map[api.ParamType]float32{
		api.SO2: 20,
		api.O3:  10,
	}
	result := mergeAverages(openMeteo, openAq)

	if result[api.SO2] != 35 {
		t.Errorf("Expected 35, got %f", result[api.SO2])
	}
	if result[api.CH4] != 30 {
		t.Errorf("Expected 30, got %f", result[api.CH4])
	}
	if result[api.O3] != 10 {
		t.Errorf("Expected 10, got %f", result[api.O3])
	}
}

func TestCalculateAverages(t *testing.T) {
	m1 := openaq.Measurement{
		Value:       10,
		ParameterId: 1,
		StationId:   2,
	}
	m2 := openaq.Measurement{
		Value:       20,
		ParameterId: 1,
		StationId:   2,
	}
	m3 := openaq.Measurement{
		Value:       30,
		ParameterId: 1,
		StationId:   2,
	}
	m4 := openaq.Measurement{
		Value:       10,
		ParameterId: 2,
		StationId:   2,
	}
	m5 := openaq.Measurement{
		Value:       20,
		ParameterId: 2,
		StationId:   2,
	}

	grouped := map[api.ParamType][]measurable{
		api.SO2: {m1, m2, m3},
		api.CH4: {m4, m5},
	}

	result := calculateAverage(grouped)
	if result[api.SO2] != 20 {
		t.Errorf("Expected 20, got %f", result[api.SO2])
	}
	if result[api.CH4] != 15 {
		t.Errorf("Expected 15, got %f", result[api.SO2])
	}
}
