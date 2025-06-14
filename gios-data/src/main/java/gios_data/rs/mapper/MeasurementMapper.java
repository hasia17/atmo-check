package gios_data.rs.mapper;

import ext.gios.api.model.GiosDataDTOLd;
import gios.data.model.MeasurementDTO;
import gios_data.domain.model.Measurement;
import gios_data.domain.model.MeasurementContext;
import org.mapstruct.Context;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;

@Mapper(componentModel = "spring")
public interface MeasurementMapper {
    @Mapping(target = "stationId", expression = "java(context.getStationId())")
    @Mapping(target = "parameterId", expression = "java(context.getParameterId())")
    @Mapping(target = "timestamp", source = "data", qualifiedByName = "mapToLocalDateTime")
    @Mapping(target = "value", source = "wartość")
    Measurement map(GiosDataDTOLd dto, @Context MeasurementContext context);

    MeasurementDTO map(Measurement dto);

    @Named("mapToLocalDateTime")
    default LocalDateTime mapToLocalDateTime(String data) {
        if (data == null) return null;
        return LocalDateTime.parse(data, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
    }

    default OffsetDateTime map(LocalDateTime dateTime) {
        return dateTime != null ? dateTime.atOffset(ZoneOffset.UTC) : null;
    }
}
