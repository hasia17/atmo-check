package data

import (
	"log/slog"
	"openaq-data/internal/api"
	"openaq-data/internal/mock"
	"openaq-data/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var initTypesLocation = types.Location{
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
	Sensors: []types.Sensor{
		{
			Parameter: types.Parameter{Id: 100},
		},
		{
			Parameter: types.Parameter{Id: 200},
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

var initTypesParameter = types.Parameter{
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

func TestStation(t *testing.T) {
	tests := []struct {
		name          string
		giveLocations []types.Location
		wantStations  []api.Station
		wantErr       error
	}{
		{
			name: "All good",
			giveLocations: []types.Location{
				initTypesLocation,
			},
			wantStations: []api.Station{
				initApiStation,
			},
			wantErr: nil,
		},
		{
			name: "Duplicate parameter",
			giveLocations: []types.Location{
				func() types.Location {
					loc := initTypesLocation
					loc.Sensors = []types.Sensor{
						{
							Parameter: types.Parameter{Id: 100},
						},
						{
							Parameter: types.Parameter{Id: 200},
						},
						{
							Parameter: types.Parameter{Id: 200},
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &slog.Logger{}
			db := &mock.Store{
				Locations: test.giveLocations,
			}
			s := NewService(db, l)

			stations, err := s.Stations(t.Context())
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantStations, stations)
		})
	}
}

func TestParameter(t *testing.T) {
	tests := []struct {
		name           string
		giveParameters []types.Parameter
		wantParameters []api.Parameter
		wantErr        error
	}{
		{
			name: "All good",
			giveParameters: []types.Parameter{
				initTypesParameter,
			},
			wantParameters: []api.Parameter{
				initApiParameter,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &slog.Logger{}
			db := &mock.Store{
				Parameters: test.giveParameters,
			}
			s := NewService(db, l)

			params, err := s.Parameters(t.Context())
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantParameters, params)
		})
	}
}
