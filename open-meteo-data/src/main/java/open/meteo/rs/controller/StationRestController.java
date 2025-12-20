package open.meteo.rs.controller;

import lombok.RequiredArgsConstructor;
import open.meteo.api.StationsApi;
import open.meteo.domain.repository.StationRepository;
import open.meteo.model.MeasurementDTO;
import open.meteo.model.StationDTO;
import open.meteo.rs.mapper.MeasurementMapper;
import open.meteo.rs.mapper.StationMapper;
import open.meteo.service.MeasurementService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequiredArgsConstructor
//@RequestMapping("/open-meteo-data-rs")
public class StationRestController implements StationsApi {

    private final StationRepository stationRepository;
    private final StationMapper stationMapper;
    private final MeasurementMapper measurementMapper;
    private final MeasurementService measurementService;

    @Override
    public ResponseEntity<List<StationDTO>> getAllStations() {
        return ResponseEntity.ok(stationMapper.map(stationRepository.findAll()));
    }

    @Override
    public ResponseEntity<List<MeasurementDTO>> getStationMeasurements(Long stationId) {
        return ResponseEntity.ok(measurementMapper.map(measurementService.getMeasurementsForStation(stationId)));

    }
}
