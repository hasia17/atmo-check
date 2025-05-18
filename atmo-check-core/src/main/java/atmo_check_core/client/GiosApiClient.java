package atmo_check_core.client;

import atmo_check_core.model.Station;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.ArrayList;
import java.util.List;

@Service
public class GiosApiClient {

    private static final String BASE_URL = "https://api.gios.gov.pl/pjp-api/rest";
    private static final String STATIONS_ENDPOINT = "/station/findAll";
    private static final String SENSORS_ENDPOINT = "/station/sensors/";
    private static final String DATA_ENDPOINT = "/data/getData/";
    private static final String AIR_QUALITY_INDEX_ENDPOINT = "/aqindex/getIndex/";

    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;

    @Autowired
    public GiosApiClient(RestTemplate restTemplate, ObjectMapper objectMapper) {
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
    }

    public List<Station> getAllStations() {
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

                JsonNode cityNode = stationNode.path("city");
                if (!cityNode.isMissingNode()) {
                    station.setCity(cityNode.path("name").asText());
                }

                station.setAddressStreet(stationNode.path("addressStreet").asText());
                stations.add(station);
            }
        }

        return stations;
    }
}
