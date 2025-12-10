package open.meteo.service;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import open.meteo.domain.model.Station;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;

@Service
@Slf4j
public class StationService {

    private final ObjectMapper objectMapper;

    public StationService() {
        this.objectMapper = new ObjectMapper();
    }

    @PostConstruct
    public void initializeStations() throws IOException {
        InputStream inputStream = getClass()
                .getResourceAsStream("/import/stations.json");

        if (inputStream == null) {
            throw new IllegalStateException("stations.json file not found");
        }

        List<Station> stations = objectMapper.readValue(
                inputStream,
                new TypeReference<List<Station>>() {}
        );

        log.info(stations.size() + " stations loaded.");
    }
}

