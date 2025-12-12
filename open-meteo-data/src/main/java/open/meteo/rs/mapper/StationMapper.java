package open.meteo.rs.mapper;

import open.meteo.domain.model.Station;
import open.meteo.model.StationDTO;
import org.mapstruct.Mapper;

import java.util.List;

@Mapper(componentModel = "spring")
public interface StationMapper {

    List<StationDTO> map(List<Station> stations);
}
