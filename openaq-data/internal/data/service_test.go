package data

import (
	"errors"
	"openaq-data/internal/api"
	"openaq-data/internal/mock"
	"openaq-data/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var errStore = errors.New("store failure")

var initModelLocation = models.Location{
	Id:       1,
	Name:     "Test Location",
	Locality: "Test Locality",
	Timezone: "UTC",
	Coordinates: struct {
		Latitude  float64 "json:\"latitude\" bson:\"latitude\""
		Longitude float64 "json:\"longitude\" bson:\"longitude\""
	}{
		Latitude:  10.0,
		Longitude: 20.0,
	},
	Sensors: []models.Sensor{
		{
			Id:        1,
			Parameter: models.Parameter{Id: 100},
		},
		{
			Id:        2,
			Parameter: models.Parameter{Id: 200},
		},
	},
}

var initApiStation = api.Station{
	Id:           1,
	Name:         "Test Location",
	Locality:     "Test Locality",
	Timezone:     "UTC",
	Latitude:     10.0,
	Longitude:    20.0,
	ParameterIds: []int32{100, 200},
}

var initModelParameter = models.Parameter{
	Id:          10,
	Name:        "Test Param",
	Units:       "test_unit",
	DisplayName: "Test Display Name",
	Description: "",
}

var initApiParameter = api.Parameter{
	Id:          10,
	Name:        "Test Param",
	Units:       "test_unit",
	DisplayName: "Test Display Name",
	Description: func() *string { //TODO: check this
		out := ""
		return &out
	}(),
}

var initModelMeasurement = models.Measurement{
	Date: struct {
		Utc   string `json:"utc" bson:"utc"`
		Local string `json:"local" bson:"local"`
	}{
		Utc:   "2009-11-17T20:34:58.651387237Z",
		Local: "2009-11-17T15:34:58.651387237-05:00",
	},
	Value: 1.2,
	Coordinates: struct {
		Latitude  float64 "json:\"latitude\" bson:\"latitude\""
		Longitude float64 "json:\"longitude\" bson:\"longitude\""
	}{
		Latitude:  10.0,
		Longitude: 20.0,
	},
	SensorId:   1,
	LocationId: 1,
}

var initApiMeasurement = api.Measurement{
	ParameterId: 100,
	StationId:   1,
	Timestamp:   time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	Value:       1.2,
}

func TestStations(t *testing.T) {
	tests := []struct {
		name          string
		giveLocations []models.Location
		giveStoreErr  error
		wantStations  []api.Station
		wantErr       error
	}{
		{
			name: "All good",
			giveLocations: []models.Location{
				initModelLocation,
			},
			wantStations: []api.Station{
				initApiStation,
			},
			wantErr: nil,
		},
		{
			name: "Duplicate parameter",
			giveLocations: []models.Location{
				func() models.Location {
					loc := initModelLocation
					loc.Sensors = []models.Sensor{
						{
							Parameter: models.Parameter{Id: 100},
						},
						{
							Parameter: models.Parameter{Id: 200},
						},
						{
							Parameter: models.Parameter{Id: 200},
						},
					}
					return loc
				}(),
			},
			wantStations: []api.Station{
				initApiStation,
			},
			wantErr: nil,
		},
		{
			name:          "No locations",
			giveLocations: nil,
			wantStations:  nil,
			wantErr:       nil,
		},
		{
			name:         "Store error is propagated",
			giveStoreErr: errStore,
			wantStations: nil,
			wantErr:      errStore,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := zap.NewNop().Sugar()
			db := &mock.Store{
				Locations:       test.giveLocations,
				GetLocationsErr: test.giveStoreErr,
			}
			s := NewService(db, l)

			stations, err := s.Stations(t.Context())
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantStations, stations)
		})
	}
}

func TestParameters(t *testing.T) {
	tests := []struct {
		name           string
		giveParameters []models.Parameter
		giveStoreErr   error
		wantParameters []api.Parameter
		wantErr        error
	}{
		{
			name: "All good",
			giveParameters: []models.Parameter{
				initModelParameter,
			},
			wantParameters: []api.Parameter{
				initApiParameter,
			},
		},
		{
			name:           "No parameters",
			giveParameters: nil,
			wantParameters: nil,
		},
		{
			name:           "Store error is propagated",
			giveStoreErr:   errStore,
			wantParameters: nil,
			wantErr:        errStore,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := zap.NewNop().Sugar()
			db := &mock.Store{
				Parameters:       test.giveParameters,
				GetParametersErr: test.giveStoreErr,
			}
			s := NewService(db, l)

			params, err := s.Parameters(t.Context())
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantParameters, params)
		})
	}
}

func TestMeasurements(t *testing.T) {
	tests := []struct {
		name                string
		giveLocations       []models.Location
		giveMeasurements    []models.Measurement
		giveLocationByIDErr error
		giveMeasurementsErr error
		wantMeasurements    []api.Measurement
		wantErr             error
	}{
		{
			name: "All good",
			giveLocations: []models.Location{
				initModelLocation,
			},
			giveMeasurements: []models.Measurement{
				initModelMeasurement,
			},
			wantMeasurements: []api.Measurement{
				initApiMeasurement,
			},
			wantErr: nil,
		},
		{
			name: "Location has no sensors",
			giveLocations: func() []models.Location {
				loc := initModelLocation
				loc.Sensors = []models.Sensor{}
				return []models.Location{
					loc,
				}
			}(),
			giveMeasurements: []models.Measurement{
				initModelMeasurement,
			},
			wantMeasurements: nil,
			wantErr:          nil,
		},
		{
			name:             "Location not found",
			giveLocations:    nil,
			giveMeasurements: nil,
			wantMeasurements: nil,
			wantErr:          ErrLocationNotFound,
		},
		{
			name: "Bad timestamp is skipped",
			giveLocations: []models.Location{
				initModelLocation,
			},
			giveMeasurements: []models.Measurement{
				func() models.Measurement {
					m := initModelMeasurement
					m.Date.Utc = "not-a-timestamp"
					return m
				}(),
			},
			wantMeasurements: nil,
			wantErr:          nil,
		},
		{
			name: "Sensor not in location is skipped",
			giveLocations: []models.Location{
				initModelLocation,
			},
			giveMeasurements: []models.Measurement{
				func() models.Measurement {
					m := initModelMeasurement
					m.SensorId = 999
					return m
				}(),
			},
			wantMeasurements: nil,
			wantErr:          nil,
		},
		{
			name: "Mix of valid and invalid returns only valid ones",
			giveLocations: []models.Location{
				initModelLocation,
			},
			giveMeasurements: []models.Measurement{
				func() models.Measurement {
					m := initModelMeasurement
					m.Date.Utc = "not-a-timestamp"
					return m
				}(),
				initModelMeasurement,
				func() models.Measurement {
					m := initModelMeasurement
					m.SensorId = 999
					return m
				}(),
				func() models.Measurement {
					m := initModelMeasurement
					m.SensorId = 2
					return m
				}(),
			},
			wantMeasurements: []api.Measurement{
				initApiMeasurement,
				func() api.Measurement {
					am := initApiMeasurement
					am.ParameterId = 200
					return am
				}(),
			},
			wantErr: nil,
		},
		{
			name:                "GetLocationByID error is propagated",
			giveLocationByIDErr: errStore,
			wantMeasurements:    nil,
			wantErr:             errStore,
		},
		{
			name: "GetMeasurementsByLocation error is propagated",
			giveLocations: []models.Location{
				initModelLocation,
			},
			giveMeasurementsErr: errStore,
			wantMeasurements:    nil,
			wantErr:             errStore,
		},
	}
	for _, tests := range tests {
		t.Run(tests.name, func(t *testing.T) {
			db := mock.Store{
				Locations:                    tests.giveLocations,
				Measurements:                 tests.giveMeasurements,
				GetLocationByIDErr:           tests.giveLocationByIDErr,
				GetMeasurementsByLocationErr: tests.giveMeasurementsErr,
			}
			l := zap.NewNop().Sugar()
			s := NewService(&db, l)

			measurements, err := s.MeasurementsForStation(t.Context(), 1)
			assert.Equal(t, tests.wantErr, err)
			assert.Equal(t, tests.wantMeasurements, measurements)
		})
	}
}
