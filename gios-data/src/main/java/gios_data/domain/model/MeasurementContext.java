package gios_data.domain.model;

public class MeasurementContext {
    private final String stationId;
    private final String parameterId;

    public MeasurementContext(String stationId, String parameterId) {
        this.stationId = stationId;
        this.parameterId = parameterId;
    }
    public String getStationId() { return stationId; }
    public String getParameterId() { return parameterId; }
}