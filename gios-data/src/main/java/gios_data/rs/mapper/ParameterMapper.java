package gios_data.rs.mapper;

import ext.gios.api.model.GiosSensorLdDTO;
import gios.data.model.ParameterDTO;
import gios_data.domain.model.Parameter;
import org.mapstruct.AfterMapping;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingTarget;

import java.util.List;

@Mapper(componentModel = "spring")
public interface ParameterMapper {
    @Mapping(target = "id", source = "idWskaźnika")
    @Mapping(target = "name", source = "wskaźnik")
    @Mapping(target = "description", source = "wskaźnikWzór")
    Parameter map(GiosSensorLdDTO dto);

    List<Parameter> map(List<GiosSensorLdDTO> dtos);

    ParameterDTO map(Parameter dto);

    @AfterMapping
    default void setUnit(GiosSensorLdDTO dto, @MappingTarget Parameter parameter) {
        parameter.setUnit(getUnitForParameterCode(dto.getWskaźnikKod()));
    }

    default String getUnitForParameterCode(String code) {
        if (code == null || code.isEmpty()) return "";
        return switch (code.toUpperCase()) {
            case "PM10", "PM2.5", "SO2", "NO2", "CO", "O3", "C6H6" -> "μg/m³";
            case "TEMP", "TEMPERATURE" -> "°C";
            case "HUMIDITY" -> "%";
            case "PRESSURE" -> "hPa";
            case "WIND_SPEED" -> "m/s";
            case "WIND_DIRECTION" -> "°";
            case "RAINFALL" -> "mm";
            default -> "";
        };
    }
}
