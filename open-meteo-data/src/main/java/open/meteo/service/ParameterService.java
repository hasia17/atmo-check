package open.meteo.service;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import lombok.AllArgsConstructor;
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
public class ParameterService {

    private final ParameterRepository parameterRepository;
    private final ObjectMapper objectMapper;

    public ParameterService(ParameterRepository parameterRepository) {
        this.parameterRepository = parameterRepository;
        this.objectMapper = new ObjectMapper();
    }


    @PostConstruct
    public void initializeParameters() throws IOException {
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
}