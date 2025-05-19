package atmo_check_core.client;

import atmo_check_core.model.Station;
import atmo_check_core.repository.CityRepository;
import atmo_check_core.repository.ParamRepository;
import atmo_check_core.repository.SensorRepository;
import atmo_check_core.repository.StationRepository;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Service
public class GiosApiClient {

    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";

    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;
    private final CityRepository cityRepository;
    private final StationRepository stationRepository;
    private final SensorRepository sensorRepository;
    private final ParamRepository paramRepository;

    @Autowired
    public GiosApiClient(
            RestTemplate restTemplate,
            ObjectMapper objectMapper,
            CityRepository cityRepository,
            StationRepository stationRepository,
            SensorRepository sensorRepository,
            ParamRepository paramRepository) {
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
        this.cityRepository = cityRepository;
        this.stationRepository = stationRepository;
        this.sensorRepository = sensorRepository;
        this.paramRepository = paramRepository;
    }

    public List<Station> fetchAllStations() {
        String url = BASE_URL + STATIONS_ENDPOINT;
        JsonNode jsonResponse = restTemplate.getForObject(url, JsonNode.class);
        List<Station> stations = new ArrayList<>();

        if (jsonResponse != null && jsonResponse.isArray()) {
            for (JsonNode stationNode : jsonResponse) {
                Station station = new Station();
                station.setId(stationNode.path("id").asInt());
                station.setStationName(stationNode.path("stationName").asText());
                station.setGegrLat(stationNode.path("gegrLat").asDouble());
                station.setGegrLon(stationNode.path("gegrLon").asDouble());
                station.setAddressStreet(stationNode.path("addressStreet").asText());

                // Sprawdź, czy stacja już istnieje w bazie danych
                Optional<Station> existingStation = stationRepository.findById(station.getId().toString());
                if (existingStation.isPresent()) {
                    // Jeśli istnieje, używamy istniejącej stacji z jej relacjami
                    station = existingStation.get();
                    // Aktualizujemy tylko podstawowe dane
                    station.setStationName(stationNode.path("stationName").asText());
                    station.setGegrLat(stationNode.path("gegrLat").asDouble());
                    station.setGegrLon(stationNode.path("gegrLon").asDouble());
                    station.setAddressStreet(stationNode.path("addressStreet").asText());
                }

                stations.add(station);
            }
        }
        return stations;
    }
}
