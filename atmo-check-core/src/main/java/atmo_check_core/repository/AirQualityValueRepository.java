package atmo_check_core.repository;

import atmo_check_core.model.AirQualityValue;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AirQualityValueRepository extends MongoRepository<AirQualityValue, Integer> {
}
