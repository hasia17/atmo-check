package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.List;

@Slf4j
@Document(collection = "sensors")
public class Sensor {
    @Id
    private Integer id;

    @DBRef
    private Param param;

    @DBRef
    private List<AirQualityValue> values;

}
