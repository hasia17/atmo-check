package aggregator.config;

import gios.data.api.MeasurementsApi;
import gios.data.api.ParametersApi;
import gios.data.api.StationsApi;
import gios.data.invoker.ApiClient;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

// This config class is needed because OPENAPI generated code requires a configuration class to instantiate the API clients. (bean not automatically created by Spring Boot)
@Configuration
public class GiosApiConfig {

    @Bean
    public ApiClient giosClient() {
        return new ApiClient();
    }

    @Bean
    public StationsApi stationsApi(ApiClient giosClient) {
        return new StationsApi(giosClient);
    }

    @Bean
    public ParametersApi parametersApi(ApiClient giosClient) {
        return new ParametersApi(giosClient);
    }

    @Bean
    public MeasurementsApi measurementsApi(ApiClient giosClient) {
        return new MeasurementsApi(giosClient);
    }
}