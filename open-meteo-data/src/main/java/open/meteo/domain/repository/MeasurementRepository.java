package open.meteo.domain.repository;

import open.meteo.domain.model.Measurement;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface MeasurementRepository extends MongoRepository<Measurement, String> {

    List<Measurement> findAllByStationId(Long stationId);

    void deleteAllByStationId(Long stationId);
}
