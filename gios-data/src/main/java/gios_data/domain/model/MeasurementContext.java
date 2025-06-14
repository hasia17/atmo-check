package gios_data.domain.model;


import lombok.Getter;

@Getter
public class MeasurementContext {
    private final String stationId;
    private final String parameterId;

    public MeasurementContext(String stationId, String parameterId) {
        this.stationId = stationId;
        this.parameterId = parameterId;
    }
}