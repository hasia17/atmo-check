package gios_data.rs.controller;

import com.example.model.SensorDTO;
import com.example.model.StationDTO;
import gios_data.domain.model.Sensor;
import gios_data.domain.model.Station;
import gios_data.domain.repository.SensorRepository;
import gios_data.domain.repository.StationRepository;
import gios_data.rs.mapper.SensorMapper;
import gios_data.rs.mapper.StationMapper;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

@Slf4j
@RestController
@RequestMapping("/gios-data")
@Tag(name = "stations", description = "Station management operations")
public class StationRestController {

    private final StationRepository stationRepository;
    private final StationMapper stationMapper;
    private final SensorRepository sensorRepository;
    private final SensorMapper sensorMapper;


    public StationRestController(StationRepository stationRepository, StationMapper stationMapper, SensorRepository sensorRepository, SensorMapper sensorMapper) {
        this.stationRepository = stationRepository;
        this.stationMapper = stationMapper;
        this.sensorRepository = sensorRepository;
        this.sensorMapper = sensorMapper;
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
    public ResponseEntity<StationDTO> getStationById(@PathVariable Integer stationId) {

        Station station = stationRepository.findById(stationId);
        return ResponseEntity.ok(stationMapper.map(station));
    }

    @GetMapping("/stations/{stationId}/sensors")
    @Operation(
            summary = "Get list of sensors for specific station",
            operationId = "getSensorsByStationId"
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "List of sensors for the station"
            ),
            @ApiResponse(
                    responseCode = "404",
                    description = "Station not found"
            )
    })
    public ResponseEntity<List<SensorDTO>> getSensorsByStationId(@PathVariable Integer stationId) {
        try {
            Station station = stationRepository.findById(stationId);
            if (station == null) {
                log.info("Station {} not found", stationId);
                return ResponseEntity.notFound().build();
            }
            List<Sensor> sensors = sensorRepository.findSensorsByStationId((stationId));

            return ResponseEntity.ok(sensorMapper.map(sensors));

        } catch (Exception e) {
            log.error("Error fetching sensors for station {}: {}", stationId, e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build();
        }
    }
}