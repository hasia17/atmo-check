package gios_data.rs.client;

import gios_data.domain.model.*;
import gios_data.domain.repository.*;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.util.*;

@Slf4j
@Service
public class GiosApiClient {

    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;
    private final StationRepository stationRepository;
    private final MeasurementRepository measurementRepository;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate,
            ObjectMapper objectMapper,
            StationRepository stationRepository,
            MeasurementRepository measurementRepository) {
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
        this.stationRepository = stationRepository;
        this.measurementRepository = measurementRepository;
    }

//    @PostConstruct
//    public void initData() {
//        updateStationsFromGios();
//    }


    @Scheduled(cron = "0 0 2 * * ?") // Daily at 2:00 AM
    public void updateStationsFromGios() {
        log.info("Starting stations update from GIOS API");

        try {
            List<Station> stations = fetchAllStations();
            log.info("Fetched {} stations from GIOS API", stations.size());

            int updatedCount = 0;
            int newCount = 0;

            for (Station station : stations) {
                try {
                    // Fetch sensors for this station0
                    List<Parameter> parameters = fetchParametersForStation(Integer.valueOf(station.getId()));

                    if (parameters.isEmpty()) {
                        log.warn("No parameters found for station: {}", station.getName());
                        continue;
                    }

                    station.setParameters(parameters);
                    station.setLastUpdated(LocalDateTime.now());

                    // Check if station exists
                    Optional<Station> existingStation = stationRepository.findById(station.getId());

                    if (existingStation.isPresent()) {
                        Station existing = existingStation.get();

                        // Check if parameters changed
                        if (!parametersEqual(existing.getParameters(), parameters)) {
                            existing.setParameters(parameters);
                            existing.setName(station.getName());
                            existing.setGegrLat(station.getGegrLat());
                            existing.setGegrLon(station.getGegrLon());
                            existing.setLastUpdated(LocalDateTime.now());

                            stationRepository.save(existing);
                            updatedCount++;
                            log.debug("Updated station: {}", station.getName());
                        }
                    } else {
                        stationRepository.save(station);
                        newCount++;
                        log.debug("Added new station: {}", station.getName());
                    }

                    // Small delay to avoid overwhelming the API
                    Thread.sleep(100);

                } catch (Exception e) {
                    log.error("Error processing station {}: {}", station.getName(), e.getMessage());
                }
            }

            log.info("Stations update completed. New: {}, Updated: {}", newCount, updatedCount);

        } catch (Exception e) {
            log.error("Failed to update stations from GIOS API", e);
            throw new RuntimeException("Station update failed", e);
        }
    }

    public List<Station> fetchAllStations() {
        log.debug("Fetching all stations from GIOS API");

        String url = BASE_URL + STATIONS_ENDPOINT;
        JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);
        List<Station> stations = new ArrayList<>();

        if (jsonResponse != null && jsonResponse.isArray()) {
            log.debug("Processing {} stations from API response", jsonResponse.size());

            for (JsonNode stationNode : jsonResponse) {
                try {
                    Station station = new Station();
                    station.setId(stationNode.path("id").asText());
                    station.setName(stationNode.path("stationName").asText());
                    station.setGegrLat(stationNode.path("gegrLat").asDouble());
                    station.setGegrLon(stationNode.path("gegrLon").asDouble());

                    stations.add(station);

                } catch (Exception e) {
                    log.error("Error parsing station data: {}", e.getMessage());
                }
            }
        } else {
            log.warn("Invalid or empty response from GIOS stations API");
        }

        log.debug("Successfully parsed {} stations", stations.size());
        return stations;
    }

    private List<Parameter> fetchParametersForStation(Integer stationId) {
        log.debug("Fetching parameters for station ID: {}", stationId);

        String url = BASE_URL + SENSORS_ENDPOINT + stationId;
        List<Parameter> parameters = new ArrayList<>();

        try {
            JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);

            if (jsonResponse != null && jsonResponse.isArray()) {
                for (JsonNode sensorNode : jsonResponse) {
                    try {
                        Parameter parameter = new Parameter();
                        parameter.setId(sensorNode.path("id").asText());

                        JsonNode paramNode = sensorNode.path("param");
                        parameter.setName(paramNode.path("paramName").asText());
                        parameter.setDescription(paramNode.path("paramFormula").asText());
                        parameter.setUnit(getUnitForParameterCode(paramNode.path("paramCode").asText()));

                        parameters.add(parameter);

                    } catch (Exception e) {
                        log.error("Error parsing sensor data for station {}: {}", stationId, e.getMessage());
                    }
                }
            }

            log.debug("Found {} parameters for station {}", parameters.size(), stationId);

        } catch (Exception e) {
            log.error("Failed to fetch parameters for station {}: {}", stationId, e.getMessage());
        }

        return parameters;
    }

    private String getUnitForParameterCode(String paramCode) {
        if (paramCode == null || paramCode.isEmpty()) {
            return "";
        }

        return switch (paramCode.toUpperCase()) {
            case "PM10", "PM2.5", "SO2", "NO2", "CO", "O3", "C6H6" -> "μg/m³";
            case "TEMP", "TEMPERATURE" -> "°C";
            case "HUMIDITY" -> "%";
            case "PRESSURE" -> "hPa";
            case "WIND_SPEED" -> "m/s";
            case "WIND_DIRECTION" -> "°";
            case "RAINFALL" -> "mm";
            default -> {
                log.debug("Unknown parameter code: {}", paramCode);
                yield "";
            }
        };
    }

    private boolean parametersEqual(List<Parameter> existing, List<Parameter> updated) {
        if (existing == null && updated == null) {
            return true;
        }
        if (existing == null || updated == null) {
            return false;
        }
        if (existing.size() != updated.size()) {
            return false;
        }

        // Sort both lists by parameterId for comparison
        List<Parameter> sortedExisting = existing.stream()
                .sorted(Comparator.comparing(Parameter::getId))
                .toList();

        List<Parameter> sortedUpdated = updated.stream()
                .sorted(Comparator.comparing(Parameter::getId))
                .toList();

        for (int i = 0; i < sortedExisting.size(); i++) {
            Parameter ex = sortedExisting.get(i);
            Parameter up = sortedUpdated.get(i);

            if (!Objects.equals(ex.getId(), up.getId()) ||
                    !Objects.equals(ex.getName(), up.getName()) ||
                    !Objects.equals(ex.getUnit(), up.getUnit()) ||
                    !Objects.equals(ex.getDescription(), up.getDescription())) {
                return false;
            }
        }

        return true;
    }
}