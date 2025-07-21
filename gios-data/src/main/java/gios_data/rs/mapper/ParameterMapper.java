package gios_data.rs.mapper;

import ext.gios.api.model.GiosSensorLdDTO;
import gios.data.model.ParameterDTO;
import gios_data.domain.model.Parameter;
import org.mapstruct.*;

import java.util.List;

@Mapper(componentModel = "spring")
public interface ParameterMapper {

    @Mapping(target = "id", source = "idWskaźnika")
    @Mapping(target = "name", source = "wskaźnik")
    @Mapping(target = "description", source = "wskaźnikWzór")
    @Mapping(target = "unit", ignore = true)
    Parameter map(GiosSensorLdDTO dto);

    List<Parameter> map(List<GiosSensorLdDTO> dtos);

    List<ParameterDTO> mapDtos(List<Parameter> dtos);

    @AfterMapping
    default void setUnit(GiosSensorLdDTO dto, @MappingTarget Parameter parameter) {
        parameter.setUnit(determineUnitForParameterCode(dto.getWskaźnikKod()));
    }

    @Named("determineUnit")
    default String determineUnitForParameterCode(String code) {
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
