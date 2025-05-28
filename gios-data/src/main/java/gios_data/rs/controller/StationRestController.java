package gios_data.rs.controller;

import com.example.model.SensorDTO;
import com.example.model.StationDTO;
import gios_data.domain.repository.StationRepository;
import gios_data.rs.mapper.StationMapper;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/gios-data")
@Tag(name = "stations", description = "Station management operations")
public class StationRestController {

    private final StationRepository stationRepository;
    private final StationMapper stationMapper;

    public StationRestController(StationRepository stationRepository, StationMapper stationMapper) {
        this.stationRepository = stationRepository;
        this.stationMapper = stationMapper;
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

        return ResponseEntity.notFound().build();
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

        return ResponseEntity.notFound().build();
    }
}