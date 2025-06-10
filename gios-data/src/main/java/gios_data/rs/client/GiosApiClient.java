package gios_data.rs.client;

import gios_data.domain.model.*;
import gios_data.domain.repository.*;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;

@Slf4j
@Service
public class GiosApiClient {

    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";
    private static final int DATA_RETENTION_DAYS = 30; // Keep measurements for 30 days
    private static final int MAX_MEASUREMENTS_PER_SENSOR = 1000; // Limit measurements per sensor
    private final RestTemplate restTemplate;
    private final StationRepository stationRepository;
    private final MeasurementRepository measurementRepository;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate,
            ObjectMapper objectMapper,
            StationRepository stationRepository,
            MeasurementRepository measurementRepository) {
        this.restTemplate = restTemplate;
        this.stationRepository = stationRepository;
        this.measurementRepository = measurementRepository;
    }


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

    @Scheduled(cron = "0 0 * * * ?") // Every hour
    public void updateMeasurementsFromGios() {
        log.info("Starting measurements update from GIOS API");

        try {
            List<Station> stations = stationRepository.findAll();
            log.info("Processing measurements for {} stations", stations.size());

            int totalMeasurements = 0;
            int newMeasurements = 0;
            int stationsProcessed = 0;

            for (Station station : stations) {
                try {
                    int stationMeasurements = updateMeasurementsForStation(station);
                    totalMeasurements += stationMeasurements;
                    if (stationMeasurements > 0) {
                        newMeasurements += stationMeasurements;
                    }
                    stationsProcessed++;

                    log.info("Measurements updated for station {}", station.getName());
                    // Small delay to avoid overwhelming the API
//                    Thread.sleep(200);

                } catch (Exception e) {
                    log.error("Error processing measurements for station {}: {}",
                            station.getName(), e.getMessage());
                }
            }

            log.info("Measurements update completed. Stations processed: {}, New measurements: {}",
                    stationsProcessed, newMeasurements);

        } catch (Exception e) {
            log.error("Failed to update measurements from GIOS API", e);
            throw new RuntimeException("Measurements update failed", e);
        }
    }

    private int updateMeasurementsForStation(Station station) {
        log.debug("Updating measurements for station: {}", station.getName());

        int newMeasurements = 0;

        for (Parameter parameter : station.getParameters()) {
            try {
                List<Measurement> measurements = fetchMeasurementsForParameter(
                        station.getId(), parameter.getId());

                for (Measurement measurement : measurements) {
                    // Check if measurement already exists
                    if (!measurementExists(measurement)) {
                        measurementRepository.save(measurement);
                        newMeasurements++;
                        log.debug("Saved new measurement for station {} parameter {} at {}",
                                station.getName(), parameter.getName(), measurement.getTimestamp());
                    }
                }

                // Small delay between parameters
//                Thread.sleep(50);

            } catch (Exception e) {
                log.error("Error fetching measurements for station {} parameter {}: {}",
                        station.getName(), parameter.getName(), e.getMessage());
            }
        }

        log.debug("Processed {} new measurements for station: {}", newMeasurements, station.getName());
        return newMeasurements;
    }

    private List<Measurement> fetchMeasurementsForParameter(String stationId, String parameterId) {
        log.debug("Fetching measurements for station {} parameter {}", stationId, parameterId);

        String url = BASE_URL + DATA_ENDPOINT + parameterId;
        List<Measurement> measurements = new ArrayList<>();

        try {
            JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);

            if (jsonResponse != null) {
                JsonNode valuesNode = jsonResponse.path("values");

                if (valuesNode.isArray()) {
                    LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);
                    int count = 0;

                    for (JsonNode valueNode : valuesNode) {
                        try {
                            // Parse date
                            String dateStr = valueNode.path("date").asText();
                            LocalDateTime timestamp = parseGiosDate(dateStr);

                            // Skip old measurements
                            if (timestamp.isBefore(cutoffDate)) {
                                continue;
                            }

                            // Parse value
                            JsonNode valueField = valueNode.path("value");
                            if (valueField.isNull() || valueField.asText().isEmpty()) {
                                continue; // Skip null/empty values
                            }

                            Double value = valueField.asDouble();
                            if (value.isNaN() || value.isInfinite()) {
                                continue; // Skip invalid values
                            }

                            Measurement measurement = new Measurement();
                            measurement.setStationId(stationId);
                            measurement.setParameterId(parameterId);
                            measurement.setValue(value);
                            measurement.setTimestamp(timestamp);

                            measurements.add(measurement);
                            count++;

                            // Limit measurements per parameter to avoid memory issues
                            if (count >= MAX_MEASUREMENTS_PER_SENSOR) {
                                log.debug("Reached max measurements limit for parameter {}", parameterId);
                                break;
                            }

                        } catch (Exception e) {
                            log.debug("Error parsing measurement value: {}", e.getMessage());
                        }
                    }

                    log.debug("Fetched {} valid measurements for parameter {}", measurements.size(), parameterId);
                }
            }

        } catch (Exception e) {
            log.error("Failed to fetch measurements for parameter {}: {}", parameterId, e.getMessage());
        }

        return measurements;
    }

    private LocalDateTime parseGiosDate(String dateStr) {
        try {
            // GIOS API uses format: "2025-06-07 14:00:00"
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
            return LocalDateTime.parse(dateStr, formatter);
        } catch (Exception e) {
            log.error("Error parsing date: {}", dateStr);
            throw new RuntimeException("Invalid date format: " + dateStr, e);
        }
    }

    private boolean measurementExists(Measurement measurement) {
        // Check if measurement with same station, parameter and timestamp already exists
        Optional<Measurement> existing = measurementRepository
                .findByStationIdAndParameterIdAndTimestamp(
                        measurement.getStationId(),
                        measurement.getParameterId(),
                        measurement.getTimestamp()
                );

        return existing.isPresent();
    }

    // Helper method to clean old measurements (optional - run weekly)
    @Scheduled(cron = "0 0 3 * * SUN") // Every Sunday at 3:00 AM
    public void cleanOldMeasurements() {
        log.info("Starting cleanup of old measurements");

        try {
            LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);

            // Find and delete old measurements
            List<Measurement> oldMeasurements = measurementRepository
                    .findByTimestampBefore(cutoffDate);

            if (!oldMeasurements.isEmpty()) {
                measurementRepository.deleteAll(oldMeasurements);
                log.info("Deleted {} old measurements older than {}",
                        oldMeasurements.size(), cutoffDate);
            } else {
                log.info("No old measurements to delete");
            }

        } catch (Exception e) {
            log.error("Error during old measurements cleanup", e);
        }
    }

    // Method to get latest measurement for station and parameter
    public Optional<Measurement> getLatestMeasurement(String stationId, String parameterId) {
        return measurementRepository
                .findFirstByStationIdAndParameterIdOrderByTimestampDesc(stationId, parameterId);
    }

    // Method to get measurements for station in time range
    public List<Measurement> getMeasurementsInRange(String stationId,
                                                    LocalDateTime from,
                                                    LocalDateTime to) {
        return measurementRepository
                .findByStationIdAndTimestampBetweenOrderByTimestampDesc(stationId, from, to);
    }
}