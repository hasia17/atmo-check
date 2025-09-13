package aggregator.service;

import aggregator.model.AggregatedVoivodeshipData;
import aggregator.model.Parameter;
import aggregator.model.SourceType;
import aggregator.model.Voivodeship;
import aggregator.rs.client.GiosApiClient;
import aggregator.rs.client.OpenAqApiClient;
import aggregator.service.wrappers.AggregatedMeasurements;
import aggregator.service.wrappers.ParameterWithMeasurements;
import gios.data.model.ParameterDTO;
import gios.data.model.StationDTO;
import gios.data.model.MeasurementDTO;

import lombok.extern.slf4j.Slf4j;
import openaq.data.model.Station;
import openaq.data.model.Measurement;
import org.springframework.stereotype.Service;
import org.springframework.util.CollectionUtils;

import java.time.OffsetDateTime;
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

        // Collect all parameters from both sources with measurements
        Map<String, List<ParameterWithMeasurements>> parameterMap = new HashMap<>();

        log.info("Processing data from gios data");
        // Process GIOS parameters with measurements
        for (StationDTO station : giosStationsInVoivodeship) {
            if (station.getParameters() != null) {
                for (ParameterDTO param : station.getParameters()) {
                    String paramKey = normalizeParameterKey(param.getDescription());

                    List<MeasurementDTO> measurements = giosApiClient.getStationMeasurements(
                            station.getId(), param.getId(), 100); // Limit to last 100 measurements

                    if (CollectionUtils.isEmpty(measurements)) {
                        continue;
                    }

                    ParameterWithMeasurements paramWithMeasurements = createParameterFromGiosData(
                            param, paramKey, measurements);
                    parameterMap.computeIfAbsent(paramKey, k -> new ArrayList<>())
                            .add(paramWithMeasurements);
                }
            }
        }

        log.info("Processing data from openaq data");
        // Process OpenAQ parameters with measurements
        for (Station station : openaqStationsInVoivodeship) {
            if (station.getParameters() != null) {

                List<Measurement> measurements = openAqApiClient.getMeasurementsByStation(
                        station.getId());
                if (CollectionUtils.isEmpty(measurements)) {
                    continue;
                }

                for (openaq.data.model.Parameter param : station.getParameters()) {
                    String paramKey = normalizeParameterKey(param.getName());

                    List<Measurement> paramMeasurements = measurements.stream()
                            .filter(m -> param.getId().equals(m.getSensorId()))
                            .collect(Collectors.toList());

                    if (CollectionUtils.isEmpty(paramMeasurements)) {
                        continue;
                    }

                    ParameterWithMeasurements paramWithMeasurements = createParameterFromOpenAqData(
                            param, paramKey, paramMeasurements);
                    parameterMap.computeIfAbsent(paramKey, k -> new ArrayList<>())
                            .add(paramWithMeasurements);
                }
            }
        }

        log.info("Aggregating parameters with measurements");
        // Aggregate parameters by creating unified Parameter objects with aggregated measurements
        List<Parameter> aggregatedParameters = parameterMap.entrySet().stream()
                .map(entry -> aggregateParameterWithMeasurements(entry.getKey(), entry.getValue()))
                .collect(Collectors.toList());

        log.info("Returning response");
        return new AggregatedVoivodeshipData()
                .voivodeship(voivodeship)
                .parameters(aggregatedParameters);
    }

    private ParameterWithMeasurements createParameterFromGiosData(ParameterDTO param, String paramKey,
                                                                  List<MeasurementDTO> measurements) {
        Parameter parameter = new Parameter();
        parameter.setName(param.getName());
        parameter.setId(paramKey);
        parameter.setUnit(param.getUnit());
        parameter.setDescription(param.getDescription());
        parameter.setSource(SourceType.GIOS);

        return new ParameterWithMeasurements(parameter, measurements, null);
    }

    private ParameterWithMeasurements createParameterFromOpenAqData(openaq.data.model.Parameter param,
                                                                    String paramKey,
                                                                    List<Measurement> measurements) {
        Parameter parameter = new Parameter();
        parameter.setName(paramKey);
        parameter.setId(param.getDisplayName());
        parameter.setUnit(param.getUnits());
        parameter.setDescription(param.getDisplayName());
        parameter.setSource(SourceType.OPEN_AQ);

        return new ParameterWithMeasurements(parameter, null, measurements);
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
     * Aggregates multiple parameter data entries with their measurements into a single Parameter object
     */
    private Parameter aggregateParameterWithMeasurements(String key, List<ParameterWithMeasurements> parameterDataList) {
        if (parameterDataList.isEmpty()) {
            return new Parameter().id(key).name(key).unit("").description("");
        }

        // Aggregate basic parameter info
        String name = parameterDataList.stream()
                .map(p -> p.parameter().getName())
                .filter(Objects::nonNull)
                .findFirst()
                .orElse(key);

        String unit = parameterDataList.stream()
                .map(p -> p.parameter().getUnit())
                .filter(Objects::nonNull)
                .filter(u -> !u.isEmpty())
                .findFirst()
                .orElse("");

        String description = parameterDataList.stream()
                .map(p -> p.parameter().getDescription())
                .filter(Objects::nonNull)
                .filter(d -> !d.isEmpty())
                .findFirst()
                .orElse("");

        Set<SourceType> sources = parameterDataList.stream()
                .map(p -> p.parameter().getSource())
                .filter(Objects::nonNull)
                .collect(Collectors.toSet());

        SourceType source;
        if (sources.size() > 1) {
            source = SourceType.ALL;
        } else {
            source = sources.stream().findFirst().orElse(null);
        }

        AggregatedMeasurements aggregatedMeasurements = aggregateMeasurements(parameterDataList);

        Parameter aggregatedParameter = new Parameter()
                .id(key)
                .name(name)
                .unit(unit)
                .description(description)
                .source(source);

        aggregatedParameter.setAverageValue(aggregatedMeasurements.averageValue());
        aggregatedParameter.setMinValue(aggregatedMeasurements.minValue());
        aggregatedParameter.setMaxValue(aggregatedMeasurements.maxValue());
        aggregatedParameter.setMeasurementCount(aggregatedMeasurements.measurementCount());
        aggregatedParameter.setLatestValue(aggregatedMeasurements.latestValue());
        aggregatedParameter.setLatestTimestamp(aggregatedMeasurements.latestTimestamp());

        return aggregatedParameter;
    }

    /**
     * Aggregates measurements from both GIOS and OpenAQ sources
     */
    private AggregatedMeasurements aggregateMeasurements(List<ParameterWithMeasurements> parameterDataList) {
        List<Double> allValues = new ArrayList<>();
        OffsetDateTime latestTimestamp = null;
        Double latestValue = null;

        for (ParameterWithMeasurements paramData : parameterDataList) {
            // Process GIOS measurements
            if (paramData.giosMeasurements() != null) {
                for (MeasurementDTO measurement : paramData.giosMeasurements()) {
                    if (measurement.getValue() != null) {
                        allValues.add(measurement.getValue());

                        // Track latest measurement
                        if (measurement.getTimestamp() != null &&
                                (latestTimestamp == null || measurement.getTimestamp().isAfter(latestTimestamp))) {
                            latestTimestamp = measurement.getTimestamp();
                            latestValue = measurement.getValue();
                        }
                    }
                }
            }

            // Process OpenAQ measurements
            if (paramData.openaqMeasurements() != null) {
                for (Measurement measurement : paramData.openaqMeasurements()) {
                    if (measurement.getValue() != null) {
                        allValues.add(measurement.getValue());

                        // Track latest measurement
                        OffsetDateTime measurementTime = convertToOffsetDateTime(measurement.getDatetime());
                        if (measurementTime != null &&
                                (latestTimestamp == null || measurementTime.isAfter(latestTimestamp))) {
                            latestTimestamp = measurementTime;
                            latestValue = measurement.getValue();
                        }
                    }
                }
            }
        }

        if (allValues.isEmpty()) {
            return new AggregatedMeasurements(null, null, null, 0, null, null);
        }

        DoubleSummaryStatistics stats = allValues.stream().mapToDouble(Double::doubleValue).summaryStatistics();

        return new AggregatedMeasurements(
                stats.getAverage(),
                stats.getMin(),
                stats.getMax(),
                (int) stats.getCount(),
                latestValue,
                latestTimestamp
        );
    }

    /**
     * Converts OpenAQ MeasurementDateTime to OffsetDateTime
     * You'll need to implement this based on your MeasurementDateTime structure
     */
    private OffsetDateTime convertToOffsetDateTime(openaq.data.model.MeasurementDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        try {
            return OffsetDateTime.parse(dateTime.toString());
        } catch (Exception e) {
            log.warn("Could not parse OpenAQ measurement datetime: {}", dateTime, e);
            return null;
        }
    }

    /**
     * Extracts parameter name from OpenAQ measurement using sensorId
     */
    private String getParameterNameFromMeasurement(Measurement measurement) {
        if (measurement == null || measurement.getSensorId() == null) {
            return "UNKNOWN";
        }

        return measurement.getSensorId().toString();
    }

}