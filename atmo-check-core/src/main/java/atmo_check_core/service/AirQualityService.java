package atmo_check_core.service;

import atmo_check_core.client.GiosApiClient;
import atmo_check_core.model.Station;
import atmo_check_core.repository.AirQualityDataRepository;
import atmo_check_core.repository.AirQualityValueRepository;
import atmo_check_core.repository.SensorRepository;
import atmo_check_core.repository.StationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class AirQualityService {

    private final GiosApiClient giosApiService;
    private final StationRepository stationRepository;
    private final SensorRepository sensorRepository;
    private final AirQualityDataRepository airQualityDataRepository;
    private final AirQualityValueRepository airQualityIndexRepository;


    // Aktualizacja stacji co 24 godziny (86400000 ms)
//    @Scheduled(fixedRate = 86400000)
    public void updateStations() {
        log.info("Stations update started");
        List<Station> stations = giosApiService.fetchAllStations();
        stationRepository.saveAll(stations);
        log.info("{} stations were updated successfully", stations.size());
    }
}
