package gios_data.domain.repository;

import gios_data.domain.model.Measurement;
import org.springframework.data.domain.Pageable;
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

    List<Measurement> findByStationIdOrderByTimestampDesc(String stationId, Pageable pageable);

    List<Measurement> findByStationIdAndParameterIdOrderByTimestampDesc(String stationId, String parameterId, Pageable pageable);
}
