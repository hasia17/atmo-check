package gios_data.domain.model;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Getter
@Setter
@Document(collection = "params")
public class Param {
    @Id
    private Integer id;
    private String paramName;
    private String paramFormula;
    private String paramCode;
}
