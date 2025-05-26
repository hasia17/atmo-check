package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Getter
@Setter
@Document(collection = "airQualityData")
public class AirQualityData {
    @Id
    private String id;
    private LocalDateTime lastUpdateDate;
}
