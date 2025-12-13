package open.meteo.service;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import open.meteo.domain.model.Measurement;
import open.meteo.domain.model.Parameter;
import open.meteo.domain.model.Station;
import open.meteo.domain.model.enums.ParameterType;
import open.meteo.domain.repository.ParameterRepository;
import open.meteo.rs.client.OpenMeteoClient;
import open.meteo.rs.dto.OpenMeteoAirQualityResponse;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@Slf4j
@AllArgsConstructor
public class MeasurementService {
    
    private final OpenMeteoClient openMeteoClient;
    private final ParameterRepository parameterRepository;
    
    private void fetchAndStoreMeasurementsForStation(Station station) {

        Map<ParameterType, Double> latestValues = fetchMeasurements(station);

        List<Parameter> parameters = parameterRepository.findAll();
        Map<ParameterType, Long> parametersTypeAndIds = parameters.stream()
                .collect(HashMap::new, (map, parameter) ->
                        map.put(ParameterType.valueOf(parameter.getName()), parameter.getId()), HashMap::putAll);

        latestValues.forEach((parameterType, value) -> {
            Long parameterId = parametersTypeAndIds.get(parameterType);
            Measurement measurement = createMeasurement(station.getId(), parameterId, value, LocalDateTime.now());




            if (parameterId != null) {
                log.info("Storing measurement for station {}, parameter {}: {}", station.getName(), parameterType, value);
                // Here you would create and save a Measurement entity using the station ID, parameterId, value, and current timestamp
                // e.g., measurementRepository.save(new Measurement(...));
            } else {
                log.warn("Parameter ID not found for parameter type: {}", parameterType);
            }

        });


    }

    private Measurement createMeasurement(Long stationId, Long parameterId, Double value, LocalDateTime dateTime) {
        Measurement measurement = new Measurement();
        measurement.setStationId(stationId);
        measurement.setParameterId(parameterId);
        measurement.setValue(value);
        measurement.setTimestamp(dateTime);
        return measurement;
    }

    private Map<ParameterType, Double> fetchMeasurements(Station station) {
        log.info("Fetching measurements for {} started", station);
        OpenMeteoAirQualityResponse measurements = openMeteoClient.getAirQuality(station.getGeoLat(), station.getGeoLon());

        Map<String, List<Double>> valuesMap = measurements.getValues();

        Map<ParameterType, Double> latestValues = new HashMap<>();

        // get only latest value for each parameter
        valuesMap.forEach((parameter, values) -> {
            if (values != null && !values.isEmpty()) {
                Double latestValue = values.getLast();
                log.info("Latest value for parameter {} at station {}: {}", parameter, station.getName(), latestValue);
                latestValues.put(ParameterType.valueOf(parameter), latestValue);
            }
        });
        return latestValues;
    }
}
