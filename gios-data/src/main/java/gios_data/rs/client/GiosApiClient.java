package gios_data.rs.client;

import ext.gios.api.model.*;
import gios_data.domain.model.*;
import gios_data.rs.mapper.MeasurementMapper;
import gios_data.rs.mapper.ParameterMapper;
import gios_data.rs.mapper.StationMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;

@Slf4j
@Service
public class GiosApiClient {

    private static final int DATA_RETENTION_DAYS = 30; // Keep measurements for 30 days
    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/v1/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";
    private static final int MAX_MEASUREMENTS_PER_SENSOR = 1; // Limit measurements per sensor
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

        String url = BASE_URL + DATA_ENDPOINT + parameterId;
        GiosCurrentDataDTO dtos = restTemplate.getForObject(url, GiosCurrentDataDTO.class);

        if (dtos == null) {
            log.warn("No measurements fetched for parameter {}", parameterId);
            return Collections.emptyList();
        }

        LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);
        MeasurementContext context = new MeasurementContext(stationId, parameterId);

        List<Measurement> measurements = dtos.getListaDanychPomiarowych().stream()
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


}