package voivodeship

import (
	"aggregator/internal/api"
	"aggregator/internal/openmeteo"
	"fmt"
	"log/slog"
)

type Aggregator struct {
	openMeteoMap    Map[openmeteo.Station]
	service         *Service
	openmeteoClient *openmeteo.Client
}

func NewAggregator(service *Service, openmeteoClient *openmeteo.Client) (*Aggregator, error) {
	openMeteoMap, err := service.GroupOpenMeteoStations()
	if err != nil {
		return nil, fmt.Errorf("failed to group open meteo stations: %w", err)
	}

	return &Aggregator{
		openMeteoMap:    openMeteoMap,
		service:         service,
		openmeteoClient: openmeteoClient,
	}, nil
}

func (a *Aggregator) AggregateOpenMeteo(v api.Voivodeship) ([]api.Parameter, error) {
	all := make([]openmeteo.Measurement, 0)
	for _, s := range a.openMeteoMap[v] {
		m, err := a.openmeteoClient.GetMeasurementForStation(s.Id)
		if err != nil {
			slog.Info("Error getting open meteo measurements: ", err)
			return nil, err
		}
		all = append(all, m...)
	}

	return calculateAverage(groupByParamId(all)), nil

}

func groupByParamId(measurements []openmeteo.Measurement) map[int][]openmeteo.Measurement {
	grouped := make(map[int][]openmeteo.Measurement)

	for _, m := range measurements {
		grouped[m.ParameterId] = append(grouped[m.ParameterId], m)
	}
	return grouped
}

func calculateAverage(grouped map[int][]openmeteo.Measurement) []api.Parameter {
	params := make([]api.Parameter, 0)
	for p, mList := range grouped {
		var sum float32 = 0.0
		for _, m := range mList {
			sum += m.Value
		}
		param := api.Parameter{
			Id:    p,
			Value: sum / float32(len(mList)),
		}
		params = append(params, param)
	}
	return params
}
