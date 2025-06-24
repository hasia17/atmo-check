package aggregator.config;

import openaq.data.api.StationsApi;
import openaq.data.api.MeasurementsApi;
import openaq.data.api.ParametersApi;
import openaq.data.invoker.ApiClient;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenAqApiConfig {

    @Bean
    public ApiClient openAqClient() {
        return new ApiClient();
    }

    @Bean
    public StationsApi openAqStationsApi(ApiClient openAqClient) {
        return new StationsApi(openAqClient);
    }

    @Bean
    public MeasurementsApi openAqMeasurementsApi(ApiClient openAqClient) {
        return new MeasurementsApi(openAqClient);
    }

    @Bean
    public ParametersApi openAqParametersApi(ApiClient openAqClient) {
        return new ParametersApi(openAqClient);
    }
}