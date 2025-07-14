package aggregator.rs;

import aggregator.api.AirQualityApi;
import aggregator.model.AggregatedVoivodeshipData;
import aggregator.model.Voivodeship;
import aggregator.service.AirQualityAggregator;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.server.ResponseStatusException;


@Slf4j
@Controller
@RequestMapping("/aggregator")
public class AirQualityRestController implements AirQualityApi {

    private final AirQualityAggregator aggregatorService;

    public AirQualityRestController(AirQualityAggregator aggregatorService) {
        this.aggregatorService = aggregatorService;
    }

    @Override
    public ResponseEntity<AggregatedVoivodeshipData> airQualityVoivodeshipGet(
            Voivodeship voivodeship) {
        try {
            aggregatorService.aggregateVoivodeshipData(voivodeship);
            return ResponseEntity.ok(aggregatorService.aggregateVoivodeshipData(voivodeship));
        } catch (IllegalArgumentException e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        }
    }
}
