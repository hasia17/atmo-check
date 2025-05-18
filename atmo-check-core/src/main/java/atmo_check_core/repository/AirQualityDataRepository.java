package atmo_check_core.repository;

import atmo_check_core.model.AirQualityData;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AirQualityDataRepository extends MongoRepository<AirQualityData, String> {
}
