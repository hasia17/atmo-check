package open.meteo.service;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import open.meteo.domain.model.Station;
import open.meteo.rs.client.OpenMeteoClient;
import open.meteo.rs.dto.OpenMeteoAirQualityResponse;
import org.springframework.stereotype.Service;

@Service
@Slf4j
@AllArgsConstructor
public class MeasurementService {
    
    private final OpenMeteoClient openMeteoClient;
    
    private void fetchAndStoreMeasurementsForStation(Station station) {
        OpenMeteoAirQualityResponse measurements = openMeteoClient.getAirQuality(station.getGeoLat(), station.getGeoLon());
    }
}
