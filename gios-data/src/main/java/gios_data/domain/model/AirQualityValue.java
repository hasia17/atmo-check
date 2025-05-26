package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Getter
@Setter
@Document(collection = "airQualityValues")
public class AirQualityValue {
    @Id
    private String id;
    private LocalDateTime date;
    private Double value;
}
