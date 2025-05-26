package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.List;

@Getter
@Setter
@Document(collection = "cities")
public class City {
    @Id
    private String id;
    private String name;
    private String communeName;
    private String districtName;
    private String provinceName;

    @DBRef
    private List<Station> stations;
}
