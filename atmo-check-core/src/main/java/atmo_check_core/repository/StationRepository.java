package atmo_check_core.repository;

import atmo_check_core.model.Station;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public abstract class StationRepository implements MongoRepository<Station, Integer> {
}
