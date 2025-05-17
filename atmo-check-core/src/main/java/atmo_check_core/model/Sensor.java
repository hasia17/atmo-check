package atmo_check_core.model;

import lombok.extern.slf4j.Slf4j;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Slf4j
@Document(collection = "sensors")
public class Sensor {
    @Id
    private Integer id;
    private Integer stationId;
    private String param;
    private String paramName;
    private String paramFormula;
    private String paramCode;
    private String idParam;

}
