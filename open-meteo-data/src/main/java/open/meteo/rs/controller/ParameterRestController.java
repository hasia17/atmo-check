package open.meteo.rs.controller;

import lombok.RequiredArgsConstructor;
import open.meteo.api.ParametersApi;
import open.meteo.domain.repository.ParameterRepository;
import open.meteo.model.ParameterDTO;
import open.meteo.rs.mapper.ParameterMapper;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequiredArgsConstructor
public class ParameterRestController implements ParametersApi {

    private final ParameterRepository parameterRepository;
    private final ParameterMapper parameterMapper;

    @Override
    public ResponseEntity<List<ParameterDTO>> getParameters() {
        return ResponseEntity.ok(parameterMapper.map(parameterRepository.findAll()));
    }
}
