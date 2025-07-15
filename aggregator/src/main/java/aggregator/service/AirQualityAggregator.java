package aggregator.service;

import aggregator.model.AggregatedVoivodeshipData;
import aggregator.model.Parameter;
import aggregator.model.Voivodeship;
import aggregator.rs.client.GiosApiClient;
import aggregator.rs.client.OpenAqApiClient;
import gios.data.model.ParameterDTO;
import gios.data.model.StationDTO;

import lombok.extern.slf4j.Slf4j;
import openaq.data.model.Station;
import org.springframework.stereotype.Service;

import java.util.*;
import java.util.stream.Collectors;

@Service
@Slf4j
public class AirQualityAggregator {

    private final GiosApiClient giosApiClient;
    private final OpenAqApiClient openAqApiClient;


    public AirQualityAggregator(GiosApiClient giosApiClient, OpenAqApiClient openAqApiClient) {
        this.giosApiClient = giosApiClient;
        this.openAqApiClient = openAqApiClient;
    }

    public AggregatedVoivodeshipData aggregateVoivodeshipData(Voivodeship voivodeship) {
        List<StationDTO> giosStations = giosApiClient.getAllStations();
        List<Station> openaqStations = openAqApiClient.getStations();

        // Filter stations within the voivodeship
        log.info("Fetching data from gios data");
        List<StationDTO> giosStationsInVoivodeship = filterGiosStationsByVoivodeship(giosStations, voivodeship);
        log.info("Fetching data from openaq data");
        List<Station> openaqStationsInVoivodeship = filterOpenaqStationsByVoivodeship(openaqStations, voivodeship);

        // Collect all parameters from both sources
        Map<String, List<Parameter>> parameterMap = new HashMap<>();

        log.info("Processing data from gios data");
        // Process GIOS parameters
        for (StationDTO station : giosStationsInVoivodeship) {
            if (station.getParameters() != null) {
                for (ParameterDTO param : station.getParameters()) {
                    String paramKey = normalizeParameterKey(param.getId());
                    Parameter parameter = createParameterFromGiosData(param, paramKey);
                    parameterMap.computeIfAbsent(paramKey, k -> new ArrayList<>()).add(parameter);
                }
            }
        }

        log.info("Processing data from openaq data");
        // Process OpenAQ parameters
        for (Station station : openaqStationsInVoivodeship) {
            if (station.getParameters() != null) {
                for (openaq.data.model.Parameter param : station.getParameters()) {
                    String paramKey = normalizeParameterKey(param.getName());
                    Parameter parameter = createParameterFromOpenAqData(param, paramKey);
                    parameterMap.computeIfAbsent(paramKey, k -> new ArrayList<>())
                            .add(new Parameter());
                }
            }
        }

        log.info("Aggregating parameters");
        // Aggregate parameters by creating unified Parameter objects
        List<Parameter> aggregatedParameters = parameterMap.entrySet().stream()
                .map(entry -> aggregateParameter(entry.getKey(), entry.getValue()))
                .collect(Collectors.toList());

        log.info("Returning response");
        return new AggregatedVoivodeshipData()
                .voivodeship(voivodeship)
                .parameters(aggregatedParameters);
    }

    private Parameter createParameterFromGiosData(ParameterDTO param, String paramKey) {
        Parameter parameter = new Parameter();
        parameter.setName(paramKey);
        parameter.setId(param.getId());
        parameter.setUnit(param.getUnit());
        parameter.setDescription(param.getDescription());
        return parameter;
    }

    private Parameter createParameterFromOpenAqData(openaq.data.model.Parameter param, String paramKey) {
        Parameter parameter = new Parameter();
        parameter.setName(paramKey);
        parameter.setId(param.getDisplayName());
        parameter.setUnit(param.getUnits());
        parameter.setDescription(param.getDisplayName());
        return parameter;
    }

    /**
     * Filters GIOS stations by voivodeship based on coordinates
     */
    private List<StationDTO> filterGiosStationsByVoivodeship(List<StationDTO> stations, Voivodeship voivodeship) {
        return stations.stream()
                .filter(station -> station.getGeoLat() != null && station.getGeoLon() != null)
                .filter(station -> VoivodeshipBounds.isInVoivodeship(
                        station.getGeoLat(), station.getGeoLon(), voivodeship)).toList();
    }

    private List<Station> filterOpenaqStationsByVoivodeship(List<Station> stations, Voivodeship voivodeship) {

        return stations.stream()
                .filter(station -> station.getCoordinates() != null)
                .filter(station -> station.getCoordinates().getLatitude() != null && station.getCoordinates().getLongitude() != null)
                .filter(station -> VoivodeshipBounds.isInVoivodeship(
                        station.getCoordinates().getLatitude(), station.getCoordinates().getLongitude(), voivodeship)).toList();
    }

    /**
     * Normalizes parameter keys for consistent aggregation
     */
    private String normalizeParameterKey(String parameterName) {
        if (parameterName == null) {
            return "UNKNOWN";
        }

        // Normalize common parameter names
        String normalized = parameterName.toUpperCase().trim();

        // Map common variations to standard names
        if (normalized.contains("PM10")) {
            return "PM10";
        } else if (normalized.contains("PM2.5") || normalized.contains("PM25")) {
            return "PM2.5";
        } else if (normalized.contains("SO2")) {
            return "SO2";
        } else if (normalized.contains("NO2")) {
            return "NO2";
        } else if (normalized.contains("CO")) {
            return "CO";
        } else if (normalized.contains("O3") || normalized.contains("OZONE")) {
            return "O3";
        }

        return normalized;
    }

    /**
     * Aggregates multiple parameter data entries into a single Parameter object
     */
    private Parameter aggregateParameter(String key, List<Parameter> parameterDataList) {
        if (parameterDataList.isEmpty()) {
            return new Parameter().id(key).name(key).unit("").description("");
        }

        // Use the first non-null values for aggregation
        Parameter primary = parameterDataList.get(0);

        String name = parameterDataList.stream()
                .map(Parameter::getName)
                .filter(Objects::nonNull)
                .findFirst()
                .orElse(key);

        String unit = parameterDataList.stream()
                .map(Parameter::getUnit)
                .filter(Objects::nonNull)
                .filter(u -> !u.isEmpty())
                .findFirst()
                .orElse("");

        String description = parameterDataList.stream()
                .map(Parameter::getDescription)
                .filter(Objects::nonNull)
                .filter(d -> !d.isEmpty())
                .findFirst()
                .orElse("");

        return new Parameter()
                .id(key)
                .name(name)
                .unit(unit)
                .description(description);
    }

}
