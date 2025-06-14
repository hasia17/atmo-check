package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class Parameter {
    private String id;
    private String name;
    private String unit;
    private String description;
}
