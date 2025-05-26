package gios_data.domain.repository;

import gios_data.domain.model.AirQualityData;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AirQualityDataRepository extends MongoRepository<AirQualityData, String> {
}
