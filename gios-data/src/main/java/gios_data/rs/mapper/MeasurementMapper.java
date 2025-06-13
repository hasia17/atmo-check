package gios_data.rs.mapper;

import gios.data.model.MeasurementDTO;
import gios_data.domain.model.Measurement;
import org.mapstruct.Mapper;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;

@Mapper(componentModel = "spring")
public interface MeasurementMapper {
    MeasurementDTO map(Measurement measurement);

    default OffsetDateTime map(LocalDateTime dateTime) {
        return dateTime != null ? dateTime.atOffset(ZoneOffset.UTC) : null;
    }
}
