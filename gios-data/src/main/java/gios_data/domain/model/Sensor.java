package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.List;

@Getter
@Setter
@Document(collection = "sensors")
public class Sensor {
    @Id
    private Integer id;

    @DBRef
    private Param param;

    @DBRef
    private List<AirQualityValue> values;

}
