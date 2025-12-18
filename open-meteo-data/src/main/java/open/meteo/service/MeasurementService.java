package open.meteo.service;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import open.meteo.domain.model.Measurement;
import open.meteo.domain.model.Parameter;
import open.meteo.domain.model.Station;
import open.meteo.domain.model.enums.ParameterType;
import open.meteo.domain.repository.MeasurementRepository;
import open.meteo.domain.repository.ParameterRepository;
import open.meteo.domain.repository.StationRepository;
import open.meteo.rs.client.OpenMeteoClient;
import open.meteo.rs.dto.OpenMeteoAirQualityResponse;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@Slf4j
@AllArgsConstructor
public class MeasurementService {
    
    private final OpenMeteoClient openMeteoClient;
    private final ParameterRepository parameterRepository;
    private final StationRepository stationRepository;
    private final MeasurementRepository measurementRepository;

    private static final String TIME_PARAM = "time";

    @Scheduled(initialDelay = 60000, fixedRate = 3600000)
    public void fetchAndStoreMeasurements() {
        log.info("Starting scheduled measurement fetch");

        List<Station> stations = stationRepository.findAll();

        for (Station station : stations) {
            try {
                fetchAndStoreMeasurementsForStation(station);
            } catch (Exception e) {
                log.error("Error fetching/storing measurements for station {}: {}", station.getName(), e.getMessage());
            }
        }
    }

    private void fetchAndStoreMeasurementsForStation(Station station) {

        List<Measurement> measurementsToCreate = new ArrayList<>();
        Map<ParameterType, Object> latestValues = fetchMeasurements(station);

        List<Parameter> parameters = parameterRepository.findAll();
        Map<ParameterType, Long> parametersTypeAndIds = parameters.stream()
                .collect(HashMap::new, (map, parameter) ->
                        map.put(ParameterType.valueOf(parameter.getName()), parameter.getId()), HashMap::putAll);

        latestValues.forEach((parameterType, value) -> {
            Long parameterId = parametersTypeAndIds.get(parameterType);
            double mappedValue = value == null ? 0 : ((Number) value).doubleValue();
            measurementsToCreate.add(createMeasurement(station.getId(), parameterId, mappedValue, LocalDateTime.now()));
        });

        measurementRepository.deleteAllByStationId(station.getId());
        measurementRepository.saveAll(measurementsToCreate);
        log.info("Stored {} measurements for station {}", measurementsToCreate.size(), station.getName());
    }

    private Measurement createMeasurement(Long stationId, Long parameterId, Double value, LocalDateTime dateTime) {
        Measurement measurement = new Measurement();
        measurement.setStationId(stationId);
        measurement.setParameterId(parameterId);
        measurement.setValue(value);
        measurement.setTimestamp(dateTime);
        return measurement;
    }

    private Map<ParameterType, Object> fetchMeasurements(Station station) {
        log.info("Fetching measurements for {} started", station.getName());
        OpenMeteoAirQualityResponse measurements = openMeteoClient.getAirQuality(station.getGeoLat(), station.getGeoLon());

        Map<String, List<Object>> valuesMap = measurements.getValues();

        Map<ParameterType, Object> latestValues = new HashMap<>();

        // get only latest value for each parameter
        valuesMap.forEach((parameter, values) -> {
            if (values != null && !values.isEmpty()) {
                if (!TIME_PARAM.equals(parameter)) {
                    Object latestValue = values.getLast();
                    log.info("Latest value for parameter {} at station {}: {}", parameter, station.getName(), latestValue);
                    latestValues.put(ParameterType.fromName(parameter), latestValue);
                }
            }
        });
        return latestValues;
    }
}
