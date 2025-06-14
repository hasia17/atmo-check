package gios_data.rs.client;

import ext.gios.api.model.GiosDataDTOLd;
import ext.gios.api.model.GiosSensorLdDTO;
import ext.gios.api.model.GiosStationLdDTO;
import gios_data.domain.model.*;
import gios_data.domain.repository.*;
import gios_data.rs.mapper.MeasurementMapper;
import gios_data.rs.mapper.ParameterMapper;
import gios_data.rs.mapper.StationMapper;
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
    private static final int MAX_MEASUREMENTS_PER_SENSOR = 1; // Limit measurements per sensor
    private final RestTemplate restTemplate;
    private final StationRepository stationRepository;
    private final MeasurementRepository measurementRepository;
    private final StationMapper stationMapper;
    private final ParameterMapper parameterMapper;
    private final MeasurementMapper measurementMapper;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate,
            StationRepository stationRepository,
            MeasurementRepository measurementRepository, StationMapper stationMapper, ParameterMapper parameterMapper, MeasurementMapper measurementMapper) {
        this.restTemplate = restTemplate;
        this.stationRepository = stationRepository;
        this.measurementRepository = measurementRepository;
        this.stationMapper = stationMapper;
        this.parameterMapper = parameterMapper;
        this.measurementMapper = measurementMapper;
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
                    List<Parameter> parameters = fetchParametersForStation(Long.valueOf(station.getId()));

                    if (parameters.isEmpty()) {
                        log.warn("No parameters found for station: {}", station.getName());
                        continue;
                    }

                    station.setParameters(parameters);
                    station.setLastUpdated(LocalDateTime.now());

                    Optional<Station> existingStation = stationRepository.findById(station.getId());

                    if (existingStation.isPresent()) {
                        Station existing = existingStation.get();

                        if (!parametersEqual(existing.getParameters(), parameters)) {
                            existing.setParameters(parameters);
                            existing.setName(station.getName());
                            existing.setGeoLat(station.getGeoLat());
                            existing.setGeoLon(station.getGeoLon());
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
        String url = BASE_URL + STATIONS_ENDPOINT;
        GiosStationLdDTO[] giosDtos = restTemplate.getForObject(url, GiosStationLdDTO[].class);

        if (giosDtos != null) {
            log.debug("Fetched {} stations from GIOS API", giosDtos.length);
            return stationMapper.mapGiosList(Arrays.asList(giosDtos));
        } else {
            log.warn("No stations fetched from GIOS API");
            return Collections.emptyList();
        }
    }

    private List<Parameter> fetchParametersForStation(Long stationId) {
        log.debug("Fetching parameters for station ID: {}", stationId);

        String url = BASE_URL + SENSORS_ENDPOINT + stationId;
        GiosSensorLdDTO[] sensors = restTemplate.getForObject(url, GiosSensorLdDTO[].class);

        if (sensors != null) {
            List<Parameter> parameters = parameterMapper.map(Arrays.asList(sensors));
            log.debug("Found {} parameters for station {}", parameters.size(), stationId);
            return parameters;
        } else {
            log.warn("No parameters fetched for station {}", stationId);
            return Collections.emptyList();
        }
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
//                     Small delay to avoid overwhelming the API
                    Thread.sleep(200);

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
                Thread.sleep(50);

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
        GiosDataDTOLd[] dtos = restTemplate.getForObject(url, GiosDataDTOLd[].class);

        if (dtos == null) {
            log.warn("No measurements fetched for parameter {}", parameterId);
            return Collections.emptyList();
        }

        LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);
        MeasurementContext context = new MeasurementContext(stationId, parameterId);

        List<Measurement> measurements = Arrays.stream(dtos)
                .filter(dto -> (dto.getData() != null)
                        && (dto.getWartość() != null)
                        && !dto.getWartość().isNaN()
                        && !dto.getWartość().isInfinite()
                        && isAfterCutoff(dto.getData(), cutoffDate)
                )
                .limit(MAX_MEASUREMENTS_PER_SENSOR)
                .map(dto -> measurementMapper.map(dto, context))
                .toList();

        log.debug("Fetched {} valid measurements for parameter {}", measurements.size(), parameterId);

        return measurements;
    }

    private boolean isAfterCutoff(String dateStr, LocalDateTime cutoff) {
        try {
            LocalDateTime ts = LocalDateTime.parse(dateStr, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
            return ts.isAfter(cutoff);
        } catch (Exception e) {
            return false;
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
}