package gios_data.domain.repository;

import gios_data.domain.model.AirQualityValue;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AirQualityValueRepository extends MongoRepository<AirQualityValue, Integer> {
}
