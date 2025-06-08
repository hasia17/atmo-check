package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.CompoundIndex;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Document(collection = "measurements")
@CompoundIndex(name = "station_param_time_idx",
        def = "{'stationId': 1, 'parameterId': 1, 'timestamp': -1}")
@CompoundIndex(name = "station_time_idx",
        def = "{'stationId': 1, 'timestamp': -1}")
@Getter
@Setter
public class Measurement {
    @Id
    private String id;
    private String stationId;
    private String parameterId;
    private Double value;
    private LocalDateTime timestamp;
}
