package open.meteo.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Getter
@Setter
@Document(collection = "stations")
public class Station {

    @Id
    private Long id;
    private String name;
    private double geoLat;
    private double geoLon;
}
