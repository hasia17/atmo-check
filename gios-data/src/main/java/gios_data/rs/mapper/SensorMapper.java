package gios_data.rs.mapper;

import com.example.model.SensorDTO;
import gios_data.domain.model.Sensor;
import org.mapstruct.Mapper;

import java.util.List;

@Mapper(componentModel = "spring")
public interface SensorMapper {

    List<SensorDTO> map(List<Sensor> sensors);
}
