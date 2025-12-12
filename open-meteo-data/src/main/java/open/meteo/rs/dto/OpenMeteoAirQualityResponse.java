package open.meteo.rs.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.util.List;

@Data
public class OpenMeteoAirQualityResponse {
    private Double latitude;
    private Double longitude;
    private Double elevation;

    @JsonProperty("generationtime_ms")
    private Double generationtimeMs;

    @JsonProperty("utc_offset_seconds")
    private Integer utcOffsetSeconds;

    private String timezone;
    private HourlyData hourly;

    @JsonProperty("hourly_units")
    private HourlyUnits hourlyUnits;
}

@Data
class HourlyData {
    private List<String> time;
    private List<Double> pm10;

    @JsonProperty("pm2_5")
    private List<Double> pm25;

    @JsonProperty("carbon_monoxide")
    private List<Double> carbonMonoxide;

    @JsonProperty("nitrogen_dioxide")
    private List<Double> nitrogenDioxide;

    @JsonProperty("sulphur_dioxide")
    private List<Double> sulphurDioxide;

    private List<Double> ozone;
}

@Data
class HourlyUnits {
    private String time;
    private String pm10;

    @JsonProperty("pm2_5")
    private String pm25;

    @JsonProperty("carbon_monoxide")
    private String carbonMonoxide;

    @JsonProperty("nitrogen_dioxide")
    private String nitrogenDioxide;

    @JsonProperty("sulphur_dioxide")
    private String sulphurDioxide;

    private String ozone;
}
