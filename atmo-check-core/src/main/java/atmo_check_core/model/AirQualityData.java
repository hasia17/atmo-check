package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Slf4j
@Document(collection = "airQualityData")
public class AirQualityData {
    @Id
    private String id;
    private Integer sensorId;
    private Double value;
    private LocalDateTime timestamp;
    private String paramCode;
    private Integer stationId;
}
