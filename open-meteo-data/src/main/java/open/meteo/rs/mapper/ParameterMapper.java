package open.meteo.rs.mapper;

import open.meteo.domain.model.Parameter;
import open.meteo.model.ParameterDTO;
import org.mapstruct.Mapper;

import java.util.List;

@Mapper(componentModel = "spring")

public interface ParameterMapper {

    List<ParameterDTO> map(List<Parameter> parameters);
}
