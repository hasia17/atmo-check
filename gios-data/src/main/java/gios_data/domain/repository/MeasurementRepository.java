package gios_data.domain.repository;

import gios_data.domain.model.Measurement;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface MeasurementRepository extends MongoRepository<Measurement, Integer> {
}
