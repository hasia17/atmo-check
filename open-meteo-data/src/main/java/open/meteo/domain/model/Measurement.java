package open.meteo.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Getter
@Setter
@Document(collection = "measurements")
public class Measurement {
    @Id
    private String id;
    private Long stationId;
    private Long parameterId;
    private Double value;
    private LocalDateTime timestamp;
}
