package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Getter
@Setter
@Document(collection = "stations")
public class Station {
    @Id
    private String id;
    private String name;
    private double geoLat;
    private double geoLon;
}
