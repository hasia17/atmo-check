package gios_data.service;


import gios_data.domain.model.Parameter;
import gios_data.domain.model.Station;
import gios_data.domain.repository.station.StationRepository;
import gios_data.rs.client.GiosApiClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Comparator;
import java.util.List;
import java.util.Objects;
import java.util.Optional;

@Slf4j
@Service
public class StationService {

    private final StationRepository stationRepository;
    private final GiosApiClient giosApiClient;

    @Autowired
    public StationService(
            StationRepository stationRepository, GiosApiClient giosApiClient) {
        this.stationRepository = stationRepository;
        this.giosApiClient = giosApiClient;
    }

    @Scheduled(cron = "0 0 2 * * ?") // Daily at 2:00 AM
    public void updateStationsFromGios() {
        log.info("Starting stations update from GIOS API");

        try {
            List<Station> stations = giosApiClient.fetchAllStations();
            log.info("Fetched {} stations from GIOS API", stations.size());

            int updatedCount = 0;
            int newCount = 0;

            for (Station station : stations) {
                try {
                    List<Parameter> parameters = giosApiClient.fetchParametersForStation(Long.valueOf(station.getId()));

                    if (parameters.isEmpty()) {
                        log.warn("No parameters found for station: {}", station.getName());
                        continue;
                    }

                    station.setParameters(parameters);
                    station.setLastUpdated(LocalDateTime.now());

                    Optional<Station> existingStation = stationRepository.findById(station.getId());

                    if (existingStation.isPresent()) {
                        Station existing = existingStation.get();

                        if (!parametersEqual(existing.getParameters(), parameters)) {
                            existing.setParameters(parameters);
                            existing.setName(station.getName());
                            existing.setGeoLat(station.getGeoLat());
                            existing.setGeoLon(station.getGeoLon());
                            existing.setLastUpdated(LocalDateTime.now());

                            stationRepository.save(existing);
                            updatedCount++;
                            log.debug("Updated station: {}", station.getName());
                        }
                    } else {
                        stationRepository.save(station);
                        newCount++;
                        log.debug("Added new station: {}", station.getName());
                    }

                    // Small delay to avoid overwhelming the API
                    Thread.sleep(100);

                } catch (Exception e) {
                    log.error("Error processing station {}: {}", station.getName(), e.getMessage());
                }
            }

            log.info("Stations update completed. New: {}, Updated: {}", newCount, updatedCount);

        } catch (Exception e) {
            log.error("Failed to update stations from GIOS API", e);
            throw new RuntimeException("Station update failed", e);
        }
    }

    private boolean parametersEqual(List<Parameter> existing, List<Parameter> updated) {
        if (existing == null && updated == null) {
            return true;
        }
        if (existing == null || updated == null) {
            return false;
        }
        if (existing.size() != updated.size()) {
            return false;
        }

        // Sort both lists by parameterId for comparison
        List<Parameter> sortedExisting = existing.stream()
                .sorted(Comparator.comparing(Parameter::getId))
                .toList();

        List<Parameter> sortedUpdated = updated.stream()
                .sorted(Comparator.comparing(Parameter::getId))
                .toList();

        for (int i = 0; i < sortedExisting.size(); i++) {
            Parameter ex = sortedExisting.get(i);
            Parameter up = sortedUpdated.get(i);

            if (!Objects.equals(ex.getId(), up.getId()) ||
                    !Objects.equals(ex.getName(), up.getName()) ||
                    !Objects.equals(ex.getUnit(), up.getUnit()) ||
                    !Objects.equals(ex.getDescription(), up.getDescription())) {
                return false;
            }
        }

        return true;
    }
}
