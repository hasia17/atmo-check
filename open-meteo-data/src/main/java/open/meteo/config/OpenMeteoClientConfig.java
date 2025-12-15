package open.meteo.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;

@Configuration
public class OpenMeteoClientConfig {

    @Bean
    public WebClient openMeteoWebClient() {
        return WebClient.builder()
                .baseUrl("https://air-quality-api.open-meteo.com/")
                .build();
    }
}
