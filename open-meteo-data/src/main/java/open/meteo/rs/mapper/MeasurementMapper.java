package open.meteo.rs.mapper;

import open.meteo.domain.model.Measurement;
import open.meteo.model.MeasurementDTO;
import org.mapstruct.Mapper;
import org.mapstruct.Named;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.List;

@Mapper(componentModel = "spring")
public interface MeasurementMapper {

    List<MeasurementDTO> map(List<Measurement> measurements);

    default OffsetDateTime map(LocalDateTime dateTime) {
        return dateTime != null ? dateTime.atOffset(ZoneOffset.UTC) : null;
    }
}
