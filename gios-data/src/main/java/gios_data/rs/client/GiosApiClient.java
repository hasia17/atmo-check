package gios_data.rs.client;

import gios_data.domain.model.*;
import gios_data.domain.repository.*;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import org.springframework.scheduling.annotation.Scheduled;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Slf4j
@Service
public class GiosApiClient {

    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";

    private static final int DATA_RETENTION_DAYS = 30;
    private static final int MAX_MEASUREMENTS_PER_SENSOR = 1000;

    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;
    private final CityRepository cityRepository;
    private final StationRepository stationRepository;
    private final SensorRepository sensorRepository;
    private final ParamRepository paramRepository;
    private final MeasurementRepository measurementRepository;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate,
            ObjectMapper objectMapper,
            CityRepository cityRepository,
            StationRepository stationRepository,
            SensorRepository sensorRepository,
            ParamRepository paramRepository,
            MeasurementRepository measurementRepository) {
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
        this.cityRepository = cityRepository;
        this.stationRepository = stationRepository;
        this.sensorRepository = sensorRepository;
        this.paramRepository = paramRepository;
        this.measurementRepository = measurementRepository;
    }

    @PostConstruct
    public void initData() {
        updateStations();
        updateAllSensors();
    }


    /**
     * Pobiera wszystkie stacje z API GIOS
     * Uruchamiane raz dziennie o 2:00
     */
    @Scheduled(cron = "0 0 2 * * *")
    public void updateStations() {
        log.info("Starting stations update from GIOS API");
        try {
            List<Station> stations = fetchAllStations();
            stationRepository.saveAll(stations);
            log.info("Updated {} stations", stations.size());
        } catch (Exception e) {
            log.error("Error during stations update", e);
        }
    }

    /**
     * Pobiera wszystkie stacje z API GIOS
     */
    public List<Station> fetchAllStations() {
        String url = BASE_URL + STATIONS_ENDPOINT;
        JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);
        List<Station> stations = new ArrayList<>();

        if (jsonResponse != null && jsonResponse.isArray()) {
            for (JsonNode stationNode : jsonResponse) {
                Station station = new Station();
                station.setId(stationNode.path("id").asInt());
                station.setStationName(stationNode.path("stationName").asText());
                station.setGegrLat(stationNode.path("gegrLat").asDouble());
                station.setGegrLon(stationNode.path("gegrLon").asDouble());
                station.setAddressStreet(stationNode.path("addressStreet").asText());

                // Sprawdź, czy stacja już istnieje w bazie danych
                Optional<Station> existingStation = stationRepository.findById(station.getId().toString());
                if (existingStation.isPresent()) {
                    // Jeśli istnieje, używamy istniejącej stacji z jej relacjami
                    station = existingStation.get();
                    // Aktualizujemy tylko podstawowe dane
                    station.setStationName(stationNode.path("stationName").asText());
                    station.setGegrLat(stationNode.path("gegrLat").asDouble());
                    station.setGegrLon(stationNode.path("gegrLon").asDouble());
                    station.setAddressStreet(stationNode.path("addressStreet").asText());
                }

                stations.add(station);
            }
        }
        return stations;
    }

    /**
     * Aktualizuje sensory dla wszystkich stacji
     * Uruchamiane raz dziennie o 3:00
     */
    @Scheduled(cron = "0 0 3 * * *")
    public void updateAllSensors() {
        log.info("Starting sensors update for all stations");
        try {
            List<Station> stations = stationRepository.findAll();
            int totalUpdated = 0;

            for (Station station : stations) {
                List<Sensor> sensors = fetchSensorsForStation(station.getId());
                if (!sensors.isEmpty()) {
                    sensorRepository.saveAll(sensors);
                    totalUpdated += sensors.size();
                }
            }

            log.info("Updated {} sensors for {} stations", totalUpdated, stations.size());
        } catch (Exception e) {
            log.error("Error during sensors update", e);
        }
    }

    /**
     * Pobiera sensory dla konkretnej stacji
     */
    public List<Sensor> fetchSensorsForStation(Integer stationId) {
        String url = BASE_URL + SENSORS_ENDPOINT + stationId;
        JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);
        List<Sensor> sensors = new ArrayList<>();

        if (jsonResponse != null && jsonResponse.isArray()) {
            for (JsonNode sensorNode : jsonResponse) {
                Sensor sensor = new Sensor();
                sensor.setId(sensorNode.path("id").asInt());
                sensor.setStationId(stationId);

                // Pobieraj dane parametru
                JsonNode paramNode = sensorNode.path("param");
                if (!paramNode.isMissingNode()) {
                    Param param = new Param();
                    param.setId(paramNode.path("idParam").asInt());
                    param.setParamName(paramNode.path("paramName").asText());
                    param.setParamFormula(paramNode.path("paramFormula").asText());
                    param.setParamCode(paramNode.path("paramCode").asText());

                    // Sprawdź czy parametr już istnieje
                    Optional<Param> existingParam = paramRepository.findById(param.getId());
                    if (existingParam.isPresent()) {
                        sensor.setParam(existingParam.get());
                    } else {
                        paramRepository.save(param);
                        sensor.setParam(param);
                    }
                }

                sensors.add(sensor);
            }
        }

        return sensors;
    }

//    /**
//     * Aktualizuje dane pomiarowe dla wszystkich sensorów
//     * Uruchamiane co godzinę o 15 minucie
//     */
//    @Scheduled(cron = "0 15 * * * *")
//    public void updateMeasurementData() {
//        log.info("Starting measurement data update");
//        try {
//            List<Sensor> sensors = sensorRepository.findAll();
//            int totalUpdated = 0;
//
//            for (Sensor sensor : sensors) {
//                List<Measurement> newMeasurements = fetchMeasurementData(sensor.getId());
//                if (!newMeasurements.isEmpty()) {
//                    // Zapisz tylko nowe pomiary (które nie istnieją już w bazie)
//                    List<Measurement> uniqueMeasurements = filterUniqueeMeasurements(newMeasurements, sensor.getId());
//                    if (!uniqueMeasurements.isEmpty()) {
//                        measurementRepository.saveAll(uniqueMeasurements);
//                        totalUpdated += uniqueMeasurements.size();
//                    }
//                }
//            }
//
//            log.info("Added {} new measurements", totalUpdated);
//
//            // Po aktualizacji usuń stare dane
//            cleanupOldMeasurements();
//
//        } catch (Exception e) {
//            log.error("Error during measurement data update", e);
//        }
//    }

//    /**
//     * Pobiera dane pomiarowe dla danego sensora
//     */
//    public List<Measurement> fetchMeasurementData(Integer sensorId) {
//        String url = BASE_URL + DATA_ENDPOINT + sensorId;
//        JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);
//        List<Measurement> measurements = new ArrayList<>();
//
//        if (jsonResponse != null && !jsonResponse.isMissingNode()) {
//            String key = jsonResponse.path("key").asText();
//            JsonNode valuesNode = jsonResponse.path("values");
//
//            if (valuesNode.isArray()) {
//                for (JsonNode valueNode : valuesNode) {
//                    String dateStr = valueNode.path("date").asText();
//                    Double value = valueNode.path("value").isNull() ? null : valueNode.path("value").asDouble();
//
//                    if (dateStr != null && !dateStr.isEmpty()) {
//                        Measurement measurement = new Measurement();
//                        measurement.setSensorId(sensorId);
//                        measurement.setParameterKey(key);
//                        measurement.setMeasurementDate(LocalDateTime.parse(dateStr, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")));
//                        measurement.setValue(value);
//                        measurement.setCreatedAt(LocalDateTime.now());
//
//                        measurements.add(measurement);
//                    }
//                }
//            }
//        }
//
//        return measurements;
//    }
//
//    /**
//     * Filtruje pomiary, zwracając tylko te, które nie istnieją już w bazie
//     */
//    private List<Measurement> filterUniqueeMeasurements(List<Measurement> measurements, Integer sensorId) {
//        List<Measurement> uniqueMeasurements = new ArrayList<>();
//
//        for (Measurement measurement : measurements) {
//            // Sprawdź czy pomiar już istnieje dla tego sensora w tej dacie
//            boolean exists = measurementRepository.existsBySensorIdAndMeasurementDate(
//                    sensorId, measurement.getMeasurementDate());
//
//            if (!exists) {
//                uniqueMeasurements.add(measurement);
//            }
//        }
//
//        return uniqueMeasurements;
//    }
//
//    /**
//     * Usuwa stare dane pomiarowe - strategia retencji danych
//     * Uruchamiane codziennie o 4:00
//     */
//    @Scheduled(cron = "0 0 4 * * *")
//    public void cleanupOldMeasurements() {
//        log.info("Starting cleanup of old measurement data");
//        try {
//            LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);
//
//            // Usuń dane starsze niż DATA_RETENTION_DAYS
//            long deletedByDate = measurementRepository.deleteByMeasurementDateBefore(cutoffDate);
//            log.info("Deleted {} measurements older than {} days", deletedByDate, DATA_RETENTION_DAYS);
//
//            // Dodatkowo, dla każdego sensora zachowaj tylko MAX_MEASUREMENTS_PER_SENSOR najnowszych pomiarów
//            List<Sensor> sensors = sensorRepository.findAll();
//            long totalDeletedByCount = 0;
//
//            for (Sensor sensor : sensors) {
//                long count = measurementRepository.countBySensorId(sensor.getId());
//                if (count > MAX_MEASUREMENTS_PER_SENSOR) {
//                    // Pobierz najstarsze pomiary do usunięcia
//                    List<Measurement> oldestMeasurements = measurementRepository
//                            .findBySensorIdOrderByMeasurementDateAsc(sensor.getId())
//                            .stream()
//                            .limit((int)(count - MAX_MEASUREMENTS_PER_SENSOR))
//                            .toList();
//
//                    measurementRepository.deleteAll(oldestMeasurements);
//                    totalDeletedByCount += oldestMeasurements.size();
//                }
//            }
//
//            if (totalDeletedByCount > 0) {
//                log.info("Additionally deleted {} measurements exceeding per-sensor limit", totalDeletedByCount);
//            }
//
//        } catch (Exception e) {
//            log.error("Error during old data cleanup", e);
//        }
//    }
//
//    /**
//     * Forces manual update of all data
//     */
//    public void forceFullUpdate() {
//        log.info("Starting full data update");
//        updateStations();
//        updateAllSensors();
//        updateMeasurementData();
//    }
//
//    /**
//     * Updates data for specific station (including sensors and measurements)
//     */
//    public void updateStationData(Integer stationId) {
//        log.info("Updating data for station {}", stationId);
//        try {
//            // Aktualizuj sensory dla stacji
//            List<Sensor> sensors = fetchSensorsForStation(stationId);
//            sensorRepository.saveAll(sensors);
//
//            // Aktualizuj pomiary dla każdego sensora tej stacji
//            for (Sensor sensor : sensors) {
//                List<Measurement> measurements = fetchMeasurementData(sensor.getId());
//                List<Measurement> uniqueMeasurements = filterUniqueeMeasurements(measurements, sensor.getId());
//                if (!uniqueMeasurements.isEmpty()) {
//                    measurementRepository.saveAll(uniqueMeasurements);
//                }
//            }
//
//            log.info("Completed data update for station {}", stationId);
//        } catch (Exception e) {
//            log.error("Error during data update for station {}", stationId, e);
//        }
//    }
}