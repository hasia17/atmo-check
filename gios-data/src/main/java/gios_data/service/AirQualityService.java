package gios_data.service;

import gios_data.rs.client.GiosApiClient;
import gios_data.domain.model.Station;
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


    public void updateStations() {
        log.info("Stations update started");
        giosApiService.updateStationsFromGios();
        log.info("stations were updated successfully");
    }
}
