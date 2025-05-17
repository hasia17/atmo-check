package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Slf4j
@Document(collection = "stations")
public class Station {

    @Id
    private Integer id;
    private String stationName;
    private Double gegrLat;
    private Double gegrLon;
    private String city;
    private String addressStreet;
    private List<Sensor> sensors;

}
