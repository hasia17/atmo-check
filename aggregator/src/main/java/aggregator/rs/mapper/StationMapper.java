package aggregator.rs.mapper;

import aggregator.model.StationWrapperDTO;
import gios.data.model.StationDTO;
import openaq.data.model.InternalStation;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface StationMapper {

    StationWrapperDTO map(StationDTO dto);

    @Mapping(target = "lacation", source = "locality")
    StationWrapperDTO map(InternalStation dto);
}
