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
    
}