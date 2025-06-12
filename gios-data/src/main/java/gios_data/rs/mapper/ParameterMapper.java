package gios_data.rs.mapper;

import com.example.model.ParameterDTO;
import gios_data.domain.model.Parameter;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface ParameterMapper {
    ParameterDTO map(Parameter parameter);
}
