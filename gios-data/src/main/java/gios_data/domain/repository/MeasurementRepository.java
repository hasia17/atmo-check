package gios_data.domain.repository;

import gios_data.domain.model.Measurement;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Repository
public interface MeasurementRepository extends MongoRepository<Measurement, Integer> {

    Optional<Measurement> findByStationIdAndParameterIdAndTimestamp(
            String stationId, String parameterId, LocalDateTime timestamp);

    List<Measurement> findByTimestampBefore(LocalDateTime cutoffDate);

    Optional<Measurement> findFirstByStationIdAndParameterIdOrderByTimestampDesc(String stationId, String parameterId);

    List<Measurement> findByStationIdAndTimestampBetweenOrderByTimestampDesc(String stationId, LocalDateTime from, LocalDateTime to);
}
