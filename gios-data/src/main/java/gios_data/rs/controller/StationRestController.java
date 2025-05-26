package gios_data.rs.controller;

import com.example.model.StationDTO;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/gios-data")
@Tag(name = "stations", description = "Station management operations")
public class StationRestController {

    @GetMapping("/stations")
    @Operation(
            summary = "Get list with all stations",
            operationId = "getAllStations"
    )
    @ApiResponses(value = {
            @ApiResponse(
                    responseCode = "200",
                    description = "All station list"
            )
    })
    public ResponseEntity<List<StationDTO>> getAllStations() {
        return ResponseEntity.ok(List.of(new StationDTO()));
    }
}
