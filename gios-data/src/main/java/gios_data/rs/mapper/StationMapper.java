package gios_data.rs.mapper;

import com.example.model.StationDTO;
import gios_data.domain.model.Station;
import org.mapstruct.Mapper;

import java.util.List;

@Mapper
public interface StationMapper {

    List<StationDTO> map(List<Station> stations);
}
