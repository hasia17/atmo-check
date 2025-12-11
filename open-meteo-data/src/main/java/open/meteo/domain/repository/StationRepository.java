package open.meteo.domain.repository;

import open.meteo.domain.model.Station;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface StationRepository extends MongoRepository<Station, String> {
}
