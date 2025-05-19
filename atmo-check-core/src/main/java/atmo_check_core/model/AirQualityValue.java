package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Slf4j
@Document(collection = "airQualityValues")
public class AirQualityValue {
    @Id
    private String id;
    private LocalDateTime date;
    private Double value;
}
