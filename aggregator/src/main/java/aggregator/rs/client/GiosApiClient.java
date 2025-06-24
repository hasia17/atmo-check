package aggregator.rs.client;

import gios.data.api.MeasurementsApi;
import gios.data.api.ParametersApi;
import gios.data.api.StationsApi;
import gios.data.model.MeasurementDTO;
import gios.data.model.ParameterDTO;
import gios.data.model.StationDTO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClientException;

import java.util.Collections;
import java.util.List;

@Slf4j
@Service
public class GiosApiClient {

    private final StationsApi stationsApi;
    private final MeasurementsApi measurementsApi;
    private final ParametersApi parametersApi;

    public GiosApiClient(StationsApi stationsApi, MeasurementsApi measurementsApi, ParametersApi parametersApi) {
        this.stationsApi = stationsApi;
        this.measurementsApi = measurementsApi;
        this.parametersApi = parametersApi;
    }

    /**
     * Returns all GIOS stations, or empty list on error.
     */
    public List<StationDTO> getAllStations() {
        try {
            return stationsApi.getAllStations();
        } catch (RestClientException ex) {
            log.error("Error fetching stations from GIOS API", ex);
            return Collections.emptyList();
        }
    }

    /**
     * Returns selected GIOS station by ID, or null on error.
     */
    public StationDTO getStationById(String stationId) {
        try {
            return stationsApi.getStationById(stationId);
        } catch (RestClientException ex) {
            log.error("Error fetching station {} from GIOS API", stationId, ex);
            return null;
        }
    }

    /**
     * Returns a list of measurements for a station and (optionally) parameter, with a result limit.
     */
    public List<MeasurementDTO> getStationMeasurements(String stationId, String parameterId, Integer limit) {
        try {
            return measurementsApi.getStationMeasurements(stationId, parameterId, limit);
        } catch (RestClientException ex) {
            log.error("Error fetching measurements for station {} from GIOS API", stationId, ex);
            return Collections.emptyList();
        }
    }

    /**
     * Returns a list of parameters available at the station with the given ID.
     */
    public List<ParameterDTO> getStationParameters(String stationId) {
        try {
            return parametersApi.getStationParameters(stationId);
        } catch (RestClientException ex) {
            log.error("Error fetching parameters for station {} from GIOS API", stationId, ex);
            return Collections.emptyList();
        }
    }
}
