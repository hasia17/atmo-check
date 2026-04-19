package aggregator

import (
	"aggregator/internal/api"
	"aggregator/internal/openaq"
	"aggregator/internal/openmeteo"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregateData(t *testing.T) {
	openMeteoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]openmeteo.Measurement{
			{ParameterId: 1, Value: 20},
		})
	}))
	defer openMeteoServer.Close()

	openAqServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]openaq.Measurement{
			{ParameterId: 1, Value: 30},
		})
	}))
	defer openAqServer.Close()

	s := &Service{
		openmeteoClient: openmeteo.NewClientWithURL(openMeteoServer.URL),
		openaqClient:    openaq.NewClientWithURL(openAqServer.URL),
		cache: cache{
			openMeteoParameters: []openmeteo.Parameter{{Id: 1, Name: "PM10"}},
			openaqParameters:    []openaq.Parameter{{Id: 1, Name: "pm10"}},
			openMeteoMap:        Map[openmeteo.Station]{api.Malopolskie: {{Id: 1}}},
			openaqMap:           Map[openaq.Station]{api.Malopolskie: {{Id: 1}}},
		},
	}

	result, err := s.AggregateData(t.Context(), api.Malopolskie)
	assert.NoError(t, err)
	assert.Equal(t, float32(25), result.Parameters[0].Value)
}

func TestAggregateDataWithCacheError(t *testing.T) {
	s := &Service{
		cache: cache{
			err: fmt.Errorf("initialization failed"),
		},
	}
	_, err := s.AggregateData(t.Context(), api.Malopolskie)
	assert.ErrorContains(t, err, "initialization failed")
}

func TestRefreshCacheWithError(t *testing.T) {
	openMeteoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stations":
			json.NewEncoder(w).Encode([]openmeteo.Station{{Id: 1, Name: "Stacja", GeoLat: 8, GeoLon: 10}})
		case "/parameters":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer openMeteoServer.Close()

	openAqServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stations":
			json.NewEncoder(w).Encode([]openaq.Station{{Id: 1, Name: "Stacja"}})
		case "/parameters":
			json.NewEncoder(w).Encode([]openaq.Parameter{{Id: 1, Name: "pm10"}})
		}
	}))
	defer openAqServer.Close()

	s := &Service{
		openmeteoClient: openmeteo.NewClientWithURL(openMeteoServer.URL),
		openaqClient:    openaq.NewClientWithURL(openAqServer.URL),
	}
	err := s.refreshCache(t.Context())
	assert.Error(t, err)
	assert.ErrorContains(t, err, "fetching openmeteo parameters")
}

func TestRefreshCache(t *testing.T) {
	openMeteoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stations":
			json.NewEncoder(w).Encode([]openmeteo.Station{{Id: 1, Name: "Stacja", GeoLat: 8, GeoLon: 10}})
		case "/parameters":
			json.NewEncoder(w).Encode([]openmeteo.Parameter{{Id: 1, Name: "PM10"}})
		}
	}))
	defer openMeteoServer.Close()

	openAqServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stations":
			json.NewEncoder(w).Encode([]openaq.Station{{Id: 1, Name: "Stacja"}})
		case "/parameters":
			json.NewEncoder(w).Encode([]openaq.Parameter{{Id: 1, Name: "pm10"}})
		}
	}))
	defer openAqServer.Close()

	s := &Service{
		openmeteoClient: openmeteo.NewClientWithURL(openMeteoServer.URL),
		openaqClient:    openaq.NewClientWithURL(openAqServer.URL),
	}
	err := s.refreshCache(t.Context())
	assert.NoError(t, err)
	assert.Len(t, s.cache.openMeteoParameters, 1)
	assert.Len(t, s.cache.openaqParameters, 1)
	assert.Len(t, s.cache.openMeteoMap, 0)
	assert.Len(t, s.cache.openaqMap, 0)
}

func TestCalculateOpenMeteoAverages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]openmeteo.Measurement{
			{ParameterId: 1, Value: 10},
			{ParameterId: 1, Value: 20},
		})
	}))
	defer server.Close()

	s := &Service{
		openmeteoClient: openmeteo.NewClientWithURL(server.URL),
	}
	parameters := []openmeteo.Parameter{
		{Id: 1, Name: "PM10"},
	}
	stations := []openmeteo.Station{
		{Id: 1},
	}
	result, err := s.calculateOpenMeteoAverages(t.Context(), parameters, stations)
	assert.NoError(t, err)
	assert.Equal(t, float32(15), result[api.PM10])
}

func TestCalculateOpenAqAverages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]openaq.Measurement{
			{ParameterId: 1, Value: 10},
			{ParameterId: 1, Value: 20},
		})
	}))
	defer server.Close()

	s := &Service{
		openaqClient: openaq.NewClientWithURL(server.URL),
	}
	parameters := []openaq.Parameter{
		{Id: 1, Name: "pm10"},
	}
	stations := []openaq.Station{
		{Id: 1},
	}
	result, err := s.calculateOpenAqAverages(t.Context(), parameters, stations)
	assert.NoError(t, err)
	assert.Equal(t, float32(15), result[api.PM10])
}

func TestBuildOpenMeteoParameterMap(t *testing.T) {
	p1 := openmeteo.Parameter{Id: 2, Name: "PM10"}
	p2 := openmeteo.Parameter{Id: 4, Name: "CARBON_MONOXIDE"}
	p3 := openmeteo.Parameter{Id: 6, Name: "UNKNOWN_PARAM"}
	result := buildOpenMeteoParameterMap([]openmeteo.Parameter{p1, p2, p3})
	assert.Equal(t, api.PM10, result[2])
	assert.Equal(t, api.CO, result[4])
	_, exists := result[6]
	assert.False(t, exists)
}

func TestBuildOpenAqParameterMap(t *testing.T) {
	p1 := openaq.Parameter{Id: 2, Name: "pm10"}
	p2 := openaq.Parameter{Id: 4, Name: "co"}
	p3 := openaq.Parameter{Id: 6, Name: "UNKNOWN_PARAM"}
	result := buildOpenAqParameterMap([]openaq.Parameter{p1, p2, p3})
	assert.Equal(t, api.PM10, result[2])
	assert.Equal(t, api.CO, result[4])
	_, exists := result[6]
	assert.False(t, exists)
}

func TestGroupStationsByVoivodeship(t *testing.T) {
	b1 := geographicalBounds{MaxLatitude: 10, MinLatitude: 5, MaxLongitude: 10, MinLongitude: 5}
	b2 := geographicalBounds{MaxLatitude: 20, MinLatitude: 11, MaxLongitude: 20, MinLongitude: 11}
	b3 := geographicalBounds{MaxLatitude: 30, MinLatitude: 21, MaxLongitude: 30, MinLongitude: 21}
	bounds := map[api.Voivodeship]geographicalBounds{
		api.Malopolskie: b1,
		api.Mazowieckie: b2,
		api.Pomorskie:   b3,
	}
	stations := []openmeteo.Station{
		{GeoLat: 8, GeoLon: 10},
		{GeoLat: 18, GeoLon: 15},
		{GeoLat: 13, GeoLon: 15},
	}
	result := groupStationsByVoivodeship(stations, bounds)
	assert.Len(t, result[api.Malopolskie], 1)
	assert.Len(t, result[api.Mazowieckie], 2)
	assert.Len(t, result[api.Pomorskie], 0)
}

func TestStationInVoivodeship(t *testing.T) {
	bounds := geographicalBounds{MaxLatitude: 10, MinLatitude: 5, MaxLongitude: 20, MinLongitude: 5}
	station := openmeteo.Station{GeoLat: 8, GeoLon: 10}
	assert.True(t, stationInVoivodeship(station, bounds))
	station.GeoLat = 20
	assert.False(t, stationInVoivodeship(station, bounds))
	station.GeoLon = 2
	assert.False(t, stationInVoivodeship(station, bounds))
}

func TestGroupByParamId(t *testing.T) {
	paramMap := map[int]api.ParamType{1: api.PM10, 2: api.SO2}
	measurements := []openmeteo.Measurement{
		{ParameterId: 1, Value: 10},
		{ParameterId: 1, Value: 20},
		{ParameterId: 2, Value: 30},
		{ParameterId: 99, Value: 40},
	}
	result := groupByParamId(measurements, paramMap)
	assert.Len(t, result[api.PM10], 2)
	assert.Len(t, result[api.SO2], 1)
	assert.Len(t, result, 2)
}

func TestMergeAverages(t *testing.T) {
	openMeteo := map[api.ParamType]float32{api.SO2: 50, api.CH4: 30}
	openAq := map[api.ParamType]float32{api.SO2: 20, api.O3: 10}
	result := mergeAverages(openMeteo, openAq)
	assert.Equal(t, float32(35), result[api.SO2])
	assert.Equal(t, float32(30), result[api.CH4])
	assert.Equal(t, float32(10), result[api.O3])
}

func TestCalculateAverages(t *testing.T) {
	m1 := openaq.Measurement{Value: 10, ParameterId: 1, StationId: 2}
	m2 := openaq.Measurement{Value: 20, ParameterId: 1, StationId: 2}
	m3 := openaq.Measurement{Value: 30, ParameterId: 1, StationId: 2}
	m4 := openaq.Measurement{Value: 10, ParameterId: 2, StationId: 2}
	m5 := openaq.Measurement{Value: 20, ParameterId: 2, StationId: 2}
	grouped := map[api.ParamType][]measurable{
		api.SO2: {m1, m2, m3},
		api.CH4: {m4, m5},
	}
	result := calculateAverage(grouped)
	assert.Equal(t, float32(20), result[api.SO2])
	assert.Equal(t, float32(15), result[api.CH4])
}
