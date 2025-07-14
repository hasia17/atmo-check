package aggregator.rs.client;

import lombok.extern.slf4j.Slf4j;
import openaq.data.api.StationsApi;
import openaq.data.api.MeasurementsApi;
import openaq.data.api.ParametersApi;
import openaq.data.model.Station;
import openaq.data.model.Measurement;
import openaq.data.model.Parameter;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClientException;

import java.util.Collections;
import java.util.List;
import java.util.Map;

@Slf4j
@Service
public class OpenAqApiClient {

    private final StationsApi stationsApi;
    private final MeasurementsApi measurementsApi;
    private final ParametersApi parametersApi;

    public OpenAqApiClient(StationsApi stationsApi, MeasurementsApi measurementsApi, ParametersApi parametersApi) {
        this.stationsApi = stationsApi;
        this.measurementsApi = measurementsApi;
        this.parametersApi = parametersApi;
    }

    /**
     * Returns a map with a list of all OpenAQ stations, or an empty map in case of error.
     */
    public List<Station> getStations() {
        try {
            return stationsApi.getStations().get(0);
        } catch (RestClientException ex) {
            log.error("Error fetching OpenAQ stations", ex);
            return Collections.emptyList();
        }
    }

    /**
     * Returns a map with the details of a station with the given ID, or an empty map in case of error.
     */
    public Map<String, Station> getStationById(Integer id) {
        try {
            return stationsApi.getStationById(id);
        } catch (RestClientException ex) {
            log.error("Error fetching OpenAQ station with ID {}", id, ex);
            return Collections.emptyMap();
        }
    }

    /**
     * Returns a map with a list of measurements for a given station ID, or an empty map in case of error.
     */
    public Map<String, List<Measurement>> getMeasurementsByStation(Integer id) {
        try {
            return measurementsApi.getMeasurementsByStation(id);
        } catch (RestClientException ex) {
            log.error("Error fetching OpenAQ measurements for station ID {}", id, ex);
            return Collections.emptyMap();
        }
    }

    /**
     * Returns a map with a list of parameters for a given station ID, or an empty map in case of error.
     */
    public Map<String, List<Parameter>> getParametersByStation(Integer id) {
        try {
            return parametersApi.getParametersByStation(id);
        } catch (RestClientException ex) {
            log.error("Error fetching OpenAQ parameters for station ID {}", id, ex);
            return Collections.emptyMap();
        }
    }
}