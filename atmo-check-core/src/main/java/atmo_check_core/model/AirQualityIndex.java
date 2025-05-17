package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;

@Slf4j
@Document(collection = "airQualityIndices")
public class AirQualityIndex {
    @Id
    private String id;
    private Integer stationId;
    private LocalDateTime timestamp;
    private String stIndexLevel;
    private Integer stIndexLevelId;
    private String stIndexStatus;
    private String stSourceDataDate;
}
