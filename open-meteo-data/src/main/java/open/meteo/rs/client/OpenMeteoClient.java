package open.meteo.rs.client;

import lombok.RequiredArgsConstructor;
import open.meteo.domain.model.enums.ParameterType;
import open.meteo.rs.dto.OpenMeteoAirQualityResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;

import java.util.Arrays;
import java.util.stream.Collectors;

@Component
@RequiredArgsConstructor
public class OpenMeteoClient {

    private final WebClient openMeteoWebClient;

    public OpenMeteoAirQualityResponse getAirQuality(double latitude, double longitude) {

        String hourlyParams = Arrays.stream(ParameterType.values())
                .map(ParameterType::getName)
                .collect(Collectors.joining(","));

        return openMeteoWebClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/v1/air-quality")
                        .queryParam("latitude", latitude)
                        .queryParam("longitude", longitude)
                        .queryParam("hourly", hourlyParams)
                        .build())
                .retrieve()
                .bodyToMono(OpenMeteoAirQualityResponse.class)
                .block();
    }


}
