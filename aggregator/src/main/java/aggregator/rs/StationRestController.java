package aggregator.rs;

import aggregator.api.StationsApi;
import aggregator.model.Station;
import aggregator.rs.client.GiosApiClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;

import java.util.Collections;
import java.util.List;

@Slf4j
@Controller
public class StationRestController implements StationsApi {

    private final GiosApiClient giosApiClient;

    public StationRestController(GiosApiClient giosApiClient) {
        this.giosApiClient = giosApiClient;
    }

    @Override
    public ResponseEntity<List<Station>> getAllStations() {
        log.info("getAllStations: {} results", giosApiClient.getAllStations().size());
        return ResponseEntity.ok(Collections.emptyList());
    }
}