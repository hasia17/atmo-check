package gios_data.rs.client;

import ext.gios.api.model.*;
import gios_data.domain.model.*;
import gios_data.rs.mapper.MeasurementMapper;
import gios_data.rs.mapper.ParameterMapper;
import gios_data.rs.mapper.StationMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.client.RestClientException;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;

@Slf4j
@Service
public class GiosApiClient {

    private static final int DATA_RETENTION_DAYS = 30;
    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/v1/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";
    private static final String ARCHIVE_DATA_ENDPOINT = "/data/getArchiveData/";
    private static final int MAX_MEASUREMENTS_PER_SENSOR = 1;

    private static final String MANUAL_STATION_ERROR_CODE = "API-ERR-100003";

    private final RestTemplate restTemplate;
    private final StationMapper stationMapper;
    private final ParameterMapper parameterMapper;
    private final MeasurementMapper measurementMapper;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate, StationMapper stationMapper, ParameterMapper parameterMapper, MeasurementMapper measurementMapper) {
        this.restTemplate = restTemplate;
        this.stationMapper = stationMapper;
        this.parameterMapper = parameterMapper;
        this.measurementMapper = measurementMapper;
    }

    public List<Station> fetchAllStations() {
        String url = BASE_URL + STATIONS_ENDPOINT;
        GiosStationLd response = restTemplate.getForObject(url, GiosStationLd.class);

        if (response != null && response.getListaStacjiPomiarowych() != null) {
            log.debug("Fetched {} stations from GIOS API", response.getListaStacjiPomiarowych().size());
            return stationMapper.mapGiosList(response.getListaStacjiPomiarowych());
        } else {
            log.warn("No stations fetched from GIOS API");
            return Collections.emptyList();
        }
    }

    public List<Parameter> fetchParametersForStation(Long stationId) {
        log.debug("Fetching parameters for station ID: {}", stationId);

        String url = BASE_URL + SENSORS_ENDPOINT + stationId;
        GiosSensorLd sensors = restTemplate.getForObject(url, GiosSensorLd.class);

        if (sensors != null) {
            List<Parameter> parameters = parameterMapper.map(sensors.getListaStanowiskPomiarowychDlaPodanejStacji());
            log.debug("Found {} parameters for station {}", parameters.size(), stationId);
            return parameters;
        } else {
            log.warn("No parameters fetched for station {}", stationId);
            return Collections.emptyList();
        }
    }

    public List<Measurement> fetchMeasurementsForParameter(String stationId, String parameterId) {
        log.debug("Fetching measurements for station {} parameter {}", stationId, parameterId);

        List<Measurement> measurements = tryFetchCurrentData(stationId, parameterId);

        if (!measurements.isEmpty()) {
            return measurements;
        }

        return fetchArchiveData(stationId, parameterId);
    }

    private List<Measurement> tryFetchCurrentData(String stationId, String parameterId) {
        try {
            String url = BASE_URL + DATA_ENDPOINT + parameterId;
            GiosCurrentDataDTO dtos = restTemplate.getForObject(url, GiosCurrentDataDTO.class);

            if (dtos == null) {
                log.debug("No current data available for parameter {}", parameterId);
                return Collections.emptyList();
            }

            if (isManualStationError(dtos)) {
                log.debug("Parameter {} is from manual station, will try archive data", parameterId);
                return Collections.emptyList();
            }

            return processMeasurements(dtos, stationId, parameterId);

        } catch (RestClientException e) {
            log.debug("Error fetching current data for parameter {}: {}", parameterId, e.getMessage());
            return Collections.emptyList();
        }
    }

    private List<Measurement> fetchArchiveData(String stationId, String parameterId) {
        try {
            log.debug("Attempting to fetch archive data for parameter {}", parameterId);
            String url = BASE_URL + ARCHIVE_DATA_ENDPOINT + parameterId;

            GiosCurrentDataDTO dtos = restTemplate.getForObject(url, GiosCurrentDataDTO.class);

            if (dtos == null) {
                log.warn("No archive data available for parameter {}", parameterId);
                return Collections.emptyList();
            }

            List<Measurement> measurements = processMeasurements(dtos, stationId, parameterId);
            log.debug("Fetched {} archive measurements for parameter {}", measurements.size(), parameterId);
            return measurements;

        } catch (RestClientException e) {
            log.warn("Error fetching archive data for parameter {}: {}", parameterId, e.getMessage());
            return Collections.emptyList();
        }
    }

    private List<Measurement> processMeasurements(GiosCurrentDataDTO dtos, String stationId, String parameterId) {
        LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);
        MeasurementContext context = new MeasurementContext(stationId, parameterId);

        return dtos.getListaDanychPomiarowych().stream()
                .filter(dto -> (dto.getData() != null)
                        && (dto.getWartość() != null)
                        && !dto.getWartość().isNaN()
                        && !dto.getWartość().isInfinite()
                        && isAfterCutoff(dto.getData(), cutoffDate)
                )
                .limit(MAX_MEASUREMENTS_PER_SENSOR)
                .map(dto -> measurementMapper.map(dto, context))
                .toList();
    }

    private boolean isManualStationError(GiosCurrentDataDTO dtos) {
        return dtos.getListaDanychPomiarowych() == null || dtos.getListaDanychPomiarowych().isEmpty();
    }

    private boolean isAfterCutoff(String dateStr, LocalDateTime cutoff) {
        try {
            LocalDateTime ts = LocalDateTime.parse(dateStr, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
            return ts.isAfter(cutoff);
        } catch (Exception e) {
            log.debug("Could not parse date: {}", dateStr);
            return false;
        }
    }
}