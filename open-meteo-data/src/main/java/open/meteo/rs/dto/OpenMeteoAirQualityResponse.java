package open.meteo.rs.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

import java.util.List;
import java.util.Map;


@Getter
@Setter
public class OpenMeteoAirQualityResponse {
    private Double latitude;
    private Double longitude;
    private Double elevation;

    @JsonProperty("generationtime_ms")
    private Double generationtimeMs;

    @JsonProperty("utc_offset_seconds")
    private Integer utcOffsetSeconds;

    private String timezone;

    @JsonProperty("timezone_abbreviation")
    private String timezoneAbbreviation;

    @JsonProperty("hourly")
    private Map<String, List<Object>> values;

    @JsonProperty("hourly_units")
    private Map<String, String> units;



//    /**
//     * Pobiera listę wszystkich dostępnych parametrów (bez 'time')
//     */
//    public List<String> getAvailableParameters() {
//        return hourly.keySet().stream()
//                .filter(key -> !key.equals("time"))
//                .sorted()
//                .toList();
//    }
//
//    /**
//     * Pobiera wartość dla danego parametru w określonym indeksie czasu
//     */
//    public Double getValue(String paramName, int timeIndex) {
//        List<Object> values = hourly.get(paramName);
//        if (values == null || timeIndex >= values.size()) {
//            return null;
//        }
//        Object value = values.get(timeIndex);
//        if (value == null) {
//            return null;
//        }
//        return value instanceof Number ? ((Number) value).doubleValue() : null;
//    }
//
//    /**
//     * Pobiera jednostkę dla danego parametru
//     */
//    public String getUnit(String paramName) {
//        return hourlyUnits != null ? hourlyUnits.get(paramName) : null;
//    }
//
//    /**
//     * Pobiera listę timestampów
//     */
//    @SuppressWarnings("unchecked")
//    public List<String> getTimeStamps() {
//        List<Object> times = hourly.get("time");
//        if (times == null) {
//            return List.of();
//        }
//        return (List<String>) (List<?>) times;
//    }
//
//    /**
//     * Pobiera liczbę punktów czasowych
//     */
//    public int getTimePointsCount() {
//        List<Object> times = hourly.get("time");
//        return times != null ? times.size() : 0;
//    }
//
//    @Override
//    public String toString() {
//        return "OpenMeteoAirQualityResponse{" +
//                "latitude=" + latitude +
//                ", longitude=" + longitude +
//                ", timezone='" + timezone + '\'' +
//                ", elevation=" + elevation +
//                ", parametersCount=" + (hourly != null ? hourly.size() - 1 : 0) +
//                '}';
//    }
}