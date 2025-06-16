package gios_data.service;

import gios_data.domain.model.Measurement;
import gios_data.domain.model.Parameter;
import gios_data.domain.model.Station;
import gios_data.domain.repository.MeasurementRepository;
import gios_data.domain.repository.station.StationRepository;
import gios_data.rs.client.GiosApiClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Slf4j
@Service
public class MeasurementService {

    private static final int DATA_RETENTION_DAYS = 30; // Keep measurements for 30 days

    private final StationRepository stationRepository;
    private final GiosApiClient giosApiClient;
    private final MeasurementRepository measurementRepository;

    @Autowired
    public MeasurementService(
            StationRepository stationRepository, GiosApiClient giosApiClient, MeasurementRepository measurementRepository) {
        this.stationRepository = stationRepository;
        this.giosApiClient = giosApiClient;
        this.measurementRepository = measurementRepository;
    }

    // Helper method to clean old measurements (optional - run weekly)
    @Scheduled(cron = "0 0 3 * * SUN") // Every Sunday at 3:00 AM
    public void cleanOldMeasurements() {
        log.info("Starting cleanup of old measurements");

        try {
            LocalDateTime cutoffDate = LocalDateTime.now().minusDays(DATA_RETENTION_DAYS);

            // Find and delete old measurements
            List<Measurement> oldMeasurements = measurementRepository
                    .findByTimestampBefore(cutoffDate);

            if (!oldMeasurements.isEmpty()) {
                measurementRepository.deleteAll(oldMeasurements);
                log.info("Deleted {} old measurements older than {}",
                        oldMeasurements.size(), cutoffDate);
            } else {
                log.info("No old measurements to delete");
            }

        } catch (Exception e) {
            log.error("Error during old measurements cleanup", e);
        }
    }

    @Scheduled(cron = "0 0 * * * ?") // Every hour
    public void updateMeasurementsFromGios() {
        log.info("Starting measurements update from GIOS API");

        try {
            List<Station> stations = stationRepository.findAll();
            log.info("Processing measurements for {} stations", stations.size());

            int totalMeasurements = 0;
            int newMeasurements = 0;
            int stationsProcessed = 0;

            for (Station station : stations) {
                try {
                    int stationMeasurements = updateMeasurementsForStation(station);
                    totalMeasurements += stationMeasurements;
                    if (stationMeasurements > 0) {
                        newMeasurements += stationMeasurements;
                    }
                    stationsProcessed++;

                    log.info("Measurements updated for station {}", station.getName());
//                     Small delay to avoid overwhelming the API
                    Thread.sleep(200);

                } catch (Exception e) {
                    log.error("Error processing measurements for station {}: {}",
                            station.getName(), e.getMessage());
                }
            }

            log.info("Measurements update completed. Stations processed: {}, New measurements: {}",
                    stationsProcessed, newMeasurements);

        } catch (Exception e) {
            log.error("Failed to update measurements from GIOS API", e);
            throw new RuntimeException("Measurements update failed", e);
        }
    }

    private int updateMeasurementsForStation(Station station) {
        log.debug("Updating measurements for station: {}", station.getName());

        int newMeasurements = 0;

        for (Parameter parameter : station.getParameters()) {
            try {
                List<Measurement> measurements = giosApiClient.fetchMeasurementsForParameter(
                        station.getId(), parameter.getId());

                for (Measurement measurement : measurements) {
                    // Check if measurement already exists
                    if (!measurementExists(measurement)) {
                        measurementRepository.save(measurement);
                        newMeasurements++;
                        log.debug("Saved new measurement for station {} parameter {} at {}",
                                station.getName(), parameter.getName(), measurement.getTimestamp());
                    }
                }

                // Small delay between parameters
                Thread.sleep(50);

            } catch (Exception e) {
                log.error("Error fetching measurements for station {} parameter {}: {}",
                        station.getName(), parameter.getName(), e.getMessage());
            }
        }

        log.debug("Processed {} new measurements for station: {}", newMeasurements, station.getName());
        return newMeasurements;
    }

    private boolean measurementExists(Measurement measurement) {
        // Check if measurement with same station, parameter and timestamp already exists
        Optional<Measurement> existing = measurementRepository
                .findByStationIdAndParameterIdAndTimestamp(
                        measurement.getStationId(),
                        measurement.getParameterId(),
                        measurement.getTimestamp()
                );

        return existing.isPresent();
    }
}
