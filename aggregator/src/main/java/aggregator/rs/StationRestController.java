package aggregator.rs;

import aggregator.api.StationsApi;
import aggregator.model.Station;
import aggregator.rs.client.GiosApiClient;
import aggregator.service.AggregatorService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;

import java.util.Collections;
import java.util.List;

@Slf4j
@Controller
public class StationRestController implements StationsApi {

    private final AggregatorService aggregatorService;

    public StationRestController(AggregatorService aggregatorService) {
        this.aggregatorService = aggregatorService;
    }

    @Override
    public ResponseEntity<List<Station>> getAllStations() {
        aggregatorService.aggregateData();
        return ResponseEntity.ok(Collections.emptyList());
    }
}