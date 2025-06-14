package gios_data.rs.mapper;

import ext.gios.api.model.GiosStationLdDTO;
import gios.data.model.StationDTO;
import gios_data.domain.model.Station;
import org.mapstruct.IterableMapping;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;

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

    @Named("mapGios")
    @Mapping(target = "id", source = "identyfikatorStacji")
    @Mapping(target = "name", source = "nazwaStacji")
    @Mapping(target = "geoLat", source = "wgS84ΦN")
    @Mapping(target = "geoLon", source = "wgS84ΛE")
    Station mapGios(GiosStationLdDTO giosStation);

    @IterableMapping(qualifiedByName = "mapGios")
    List<Station> mapGiosList(List<GiosStationLdDTO> giosStations);
}
