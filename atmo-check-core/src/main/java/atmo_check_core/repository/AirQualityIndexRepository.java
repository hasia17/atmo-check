package atmo_check_core.repository;

import atmo_check_core.model.AirQualityIndex;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AirQualityIndexRepository extends MongoRepository<AirQualityIndex, Integer> {
}
