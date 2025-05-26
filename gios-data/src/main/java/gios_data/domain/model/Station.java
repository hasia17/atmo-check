package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.List;

@Getter
@Setter
@Document(collection = "stations")
public class Station {

    @Id
    private Integer id;
    private String stationName;
    private Double gegrLat;
    private Double gegrLon;
    private String addressStreet;

    @DBRef
    private List<Sensor> sensors;

    @DBRef
    private AirQualityData airQualityData;

}
