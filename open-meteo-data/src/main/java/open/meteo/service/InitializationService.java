package open.meteo.service;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import open.meteo.domain.model.Parameter;
import open.meteo.domain.model.Station;
import open.meteo.domain.repository.ParameterRepository;
import open.meteo.domain.repository.StationRepository;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;

@Service
@Slf4j
public class InitializationService {

    private final ParameterRepository parameterRepository;
    private final ObjectMapper objectMapper;
    private final StationRepository stationRepository;
    private final MeasurementService measurementService;


    public InitializationService(ParameterRepository parameterRepository, StationRepository stationRepository, MeasurementService measurementService) {
        this.parameterRepository = parameterRepository;
        this.stationRepository = stationRepository;
        this.measurementService = measurementService;
        this.objectMapper = new ObjectMapper();
    }

    @PostConstruct
    public void initializeData() {
        try {
            initializeParameters();
            initializeStations();
            measurementService.fetchAndStoreMeasurements();
        } catch (IOException e) {
            log.error("Error during initialization: {}", e.getMessage());
        }
    }


    private void initializeParameters() throws IOException {
        // clean database before loading
        parameterRepository.deleteAll();

        // load parameters from JSON file
        InputStream inputStream = getClass()
                .getResourceAsStream("/import/parameters.json");

        if (inputStream == null) {
            throw new IllegalStateException("parameters.json file not found");
        }

        List<Parameter> parameters = objectMapper.readValue(
                inputStream,
                new TypeReference<List<Parameter>>() {
                }
        );

        // save parameters to database
        parameterRepository.saveAll(parameters);
        log.info("{} parameters loaded.", parameters.size());
    }


    private void initializeStations() throws IOException {
        // clean database before loading
        stationRepository.deleteAll();

        // load stations from JSON file
        InputStream inputStream = getClass()
                .getResourceAsStream("/import/stations.json");

        if (inputStream == null) {
            throw new IllegalStateException("stations.json file not found");
        }

        List<Station> stations = objectMapper.readValue(
                inputStream,
                new TypeReference<List<Station>>() {}
        );

        // save stations to database
        stationRepository.saveAll(stations);
        log.info("{} stations loaded.", stations.size());
    }
}
