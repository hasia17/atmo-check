package gios_data.rs.controller;

import gios.data.model.MeasurementDTO;
import gios.data.model.ParameterDTO;
import gios.data.model.StationDTO;
import gios_data.domain.model.Measurement;
import gios_data.domain.model.Station;
import gios_data.domain.repository.MeasurementRepository;
import gios_data.domain.repository.station.StationRepository;
import gios_data.rs.mapper.MeasurementMapper;
import gios_data.rs.mapper.ParameterMapper;
import gios_data.rs.mapper.StationMapper;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;
import java.util.Optional;

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
            summary = "Get list with stations",
            operationId = "getStations"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "All station list or filtered stations"),
            @ApiResponse(responseCode = "400", description = "Bad request - invalid parameters"),
            @ApiResponse(responseCode = "500", description = "Server error")
    })
    public ResponseEntity<List<StationDTO>> getStations(
            @RequestParam(required = false) String city,
            @RequestParam(required = false) String province,
            @RequestParam(required = false) Double lat,
            @RequestParam(required = false) Double lon,
            @RequestParam(required = false, defaultValue = "10") Double radius) {

        try {
            return ResponseEntity.ok(stationMapper.map(stationRepository.searchByCriteria(city, province, lat, lon, radius)));
        } catch (IllegalArgumentException e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        }
    }

    @GetMapping("/stations/{stationId}")
    @Operation(
            summary = "Get station by ID",
            operationId = "getStationById"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Station details"),
            @ApiResponse(responseCode = "404", description = "Station not found"),
            @ApiResponse(responseCode = "500", description = "Server error")
    })
    public ResponseEntity<StationDTO> getStationById(@PathVariable String stationId) {
        try {
            Optional<Station> station = stationRepository.findById(stationId);
            return station.map(value -> ResponseEntity.ok(stationMapper.map(value))).orElseGet(() -> ResponseEntity.notFound().build());
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Server error", e);
        }
    }

    @GetMapping("/stations/{stationId}/parameters")
    @Operation(
            summary = "Get parameters for a station",
            operationId = "getStationParameters"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "List of parameters for the station"),
            @ApiResponse(responseCode = "404", description = "Station not found"),
            @ApiResponse(responseCode = "500", description = "Server error")
    })
    public ResponseEntity<List<ParameterDTO>> getParametersForStation(@PathVariable String stationId) {
        try {
            return stationRepository.findById(stationId)
                    .map(station -> ResponseEntity.ok(parameterMapper.mapDtos(station.getParameters())))
                    .orElseGet(() -> ResponseEntity.notFound().build());
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Server error", e);
        }
    }


    @GetMapping("/stations/{stationId}/measurements")
    @Operation(
            summary = "Get measurement data for a station",
            operationId = "getStationMeasurements"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Measurement data for station"),
            @ApiResponse(responseCode = "404", description = "Station or measurements not found"),
            @ApiResponse(responseCode = "500", description = "Server error")
    })
    public ResponseEntity<List<MeasurementDTO>> getStationMeasurements(
            @PathVariable String stationId,
            @RequestParam(required = false) String parameterId,
            @RequestParam(required = false, defaultValue = "100") Integer limit
    ) {
        try {
            if (!stationRepository.existsById(stationId)) {
                return ResponseEntity.notFound().build();
            }

            List<Measurement> measurements;
            Pageable pageable = PageRequest.of(0, limit);

            if (parameterId != null && !parameterId.isEmpty()) {
                measurements = measurementRepository.findByStationIdAndParameterIdOrderByTimestampDesc(stationId, parameterId, pageable);
            } else {
                measurements = measurementRepository.findByStationIdOrderByTimestampDesc(stationId, pageable);
            }

            return ResponseEntity.ok(measurementMapper.map(measurements));
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Server error", e);
        }
    }
}