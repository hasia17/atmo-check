package aggregator.rs;

import aggregator.api.StationsApi;
import aggregator.model.Station;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;

import java.util.Collections;
import java.util.List;

@Controller
public class StationRestController implements StationsApi {

    @Override
    public ResponseEntity<List<Station>> getAllStations() {
        return ResponseEntity.ok(Collections.emptyList());
    }
}