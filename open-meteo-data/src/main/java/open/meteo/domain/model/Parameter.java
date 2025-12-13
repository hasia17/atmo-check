package open.meteo.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.mongodb.core.mapping.Document;

@Getter
@Setter
@Document(collection = "parameters")
public class Parameter {
    private Long id;
    private String name;
    private String unit;
    private String description;
}
