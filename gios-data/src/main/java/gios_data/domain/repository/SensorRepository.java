package gios_data.domain.repository;

import gios_data.domain.model.Sensor;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface SensorRepository extends MongoRepository<Sensor, String> {
    List<Sensor> findSensorsByStationId(Integer stationId);
}
