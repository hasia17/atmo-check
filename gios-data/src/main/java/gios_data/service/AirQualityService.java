package gios_data.service;

import gios_data.rs.client.GiosApiClient;
import gios_data.domain.model.Station;
import gios_data.domain.repository.AirQualityDataRepository;
import gios_data.domain.repository.AirQualityValueRepository;
import gios_data.domain.repository.SensorRepository;
import gios_data.domain.repository.StationRepository;
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


    public void updateStations() {
        log.info("Stations update started");
        List<Station> stations = giosApiService.fetchAllStations();
        stationRepository.saveAll(stations);
        log.info("{} stations were updated successfully", stations.size());
    }
}
