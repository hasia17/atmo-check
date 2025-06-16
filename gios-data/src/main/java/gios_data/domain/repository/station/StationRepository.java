package gios_data.domain.repository.station;

import gios_data.domain.model.Station;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;


@Repository
public interface StationRepository extends MongoRepository<Station, String>, StationRepositoryCustom {
}

