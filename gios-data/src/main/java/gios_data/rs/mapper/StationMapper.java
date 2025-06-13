package gios_data.rs.mapper;

import gios.data.model.StationDTO;
import gios_data.domain.model.Station;
import org.mapstruct.Mapper;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.List;

@Mapper(componentModel = "spring")
public interface StationMapper {

    List<StationDTO> map(List<Station> stations);
    StationDTO map(Station station);

    default OffsetDateTime map(LocalDateTime dateTime) {
        return dateTime != null ? dateTime.atOffset(ZoneOffset.UTC) : null;
    }
}
