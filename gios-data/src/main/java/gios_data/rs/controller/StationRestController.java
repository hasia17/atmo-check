package gios_data.rs.controller;

import com.example.model.MeasurementDTO;
import com.example.model.ParameterDTO;
import com.example.model.StationDTO;
import gios_data.domain.model.Measurement;
import gios_data.domain.model.Parameter;
import gios_data.domain.model.Station;
import gios_data.domain.repository.MeasurementRepository;
import gios_data.domain.repository.StationRepository;
import gios_data.rs.mapper.MeasurementMapper;
import gios_data.rs.mapper.ParameterMapper;
import gios_data.rs.mapper.StationMapper;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.extern.slf4j.Slf4j;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.Comparator;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Slf4j
@RestController
@RequestMapping("/gios-data")
@Tag(name = "stations", description = "Station management operations")
public class StationRestController {

    private final StationRepository stationRepository;
    private final StationMapper stationMapper;
    private final MeasurementRepository measurementRepository;
    private final ParameterMapper parameterMapper;
    private final MeasurementMapper measurementMapper;


    public StationRestController(StationRepository stationRepository, StationMapper stationMapper, MeasurementRepository measurementRepository, ParameterMapper parameterMapper, MeasurementMapper measurementMapper) {
        this.stationRepository = stationRepository;
        this.stationMapper = stationMapper;
        this.measurementRepository = measurementRepository;
        this.parameterMapper = parameterMapper;
        this.measurementMapper = measurementMapper;
    }

    @GetMapping("/stations")
    @Operation(
            summary = "Get list with all stations",
            operationId = "getAllStations"
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "All station list or filtered stations"
            ),
            @ApiResponse(
                    responseCode = "400",
                    description = "Bad request - invalid parameters"
            )
    })
    public ResponseEntity<List<StationDTO>> getAllStations(
            @RequestParam(required = false) String city,
            @RequestParam(required = false) String province,
            @RequestParam(required = false) Double lat,
            @RequestParam(required = false) Double lon,
            @RequestParam(required = false, defaultValue = "10") Double radius) {

        return ResponseEntity.ok(stationMapper.map(stationRepository.findAll()));
    }

    @GetMapping("/stations/{stationId}")
    @Operation(
            summary = "Get station by ID",
            operationId = "getStationById"
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "Station details"
            ),
            @ApiResponse(
                    responseCode = "404",
                    description = "Station not found"
            )
    })
    public ResponseEntity<StationDTO> getStationById(@PathVariable String stationId) {

        Optional<Station> station = stationRepository.findById(stationId);
        return station.map(value -> ResponseEntity.ok(stationMapper.map(value))).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @GetMapping("/stations/{stationId}/parameters")
    @Operation(
            summary = "Get parameters for a station",
            operationId = "getStationParameters"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "List of parameters for the station"),
            @ApiResponse(responseCode = "404", description = "Station not found")
    })
    public ResponseEntity<List<ParameterDTO>> getParametersForStation(@PathVariable String stationId) {
        return stationRepository.findById(stationId)
                .map(station -> {
                    List<ParameterDTO> parameterDTOs = station.getParameters().stream()
                            .map(parameterMapper::map)
                            .collect(Collectors.toList());
                    return ResponseEntity.ok(parameterDTOs);
                })
                .orElseGet(() -> ResponseEntity.notFound().build());
    }

    @GetMapping("/stations/{stationId}/measurements")
    @Operation(
            summary = "Get measurement data for a station",
            operationId = "getStationMeasurements"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Measurement data for station"),
            @ApiResponse(responseCode = "404", description = "Station or measurements not found")
    })
    public ResponseEntity<List<MeasurementDTO>> getStationMeasurements(
            @PathVariable String stationId,
            @RequestParam(required = false) String parameterId,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime from,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime to,
            @RequestParam(required = false, defaultValue = "100") Integer limit
    ) {
        if (!stationRepository.existsById(stationId)) {
            return ResponseEntity.notFound().build();
        }

        List<Measurement> measurements = measurementRepository.findAll().stream()
                .filter(m -> m.getStationId().equals(stationId))
                .filter(m -> parameterId == null || parameterId.equals(m.getParameterId()))
                .filter(m -> from == null || !m.getTimestamp().isBefore(from))
                .filter(m -> to == null || !m.getTimestamp().isAfter(to))
                .sorted(Comparator.comparing(Measurement::getTimestamp).reversed())
                .limit(limit)
                .toList();

        return ResponseEntity.ok(
                measurements.stream().map(measurementMapper::map).collect(Collectors.toList())
        );
    }
}