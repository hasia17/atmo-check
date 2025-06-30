package aggregator.service;

import aggregator.rs.client.GiosApiClient;
import aggregator.rs.client.OpenAqApiClient;
import gios.data.model.StationDTO;
import openaq.data.model.InternalStation;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;

@Service
public class AggregatorService {

    private final GiosApiClient giosApiClient;
    private final OpenAqApiClient openAqApiClient;


    public AggregatorService(GiosApiClient giosApiClient, OpenAqApiClient openAqApiClient) {
        this.giosApiClient = giosApiClient;
        this.openAqApiClient = openAqApiClient;
    }

    public void aggregateData() {
        List<StationDTO> giosStations = giosApiClient.getAllStations();
        Map<String, List<InternalStation>> openAqStations = openAqApiClient.getStations();

        System.out.println("GIOS Stations: " + giosStations);
        System.out.println("OpenAQ Stations: " + openAqStations);
    }


}
