package aggregator.service.wrappers;

import aggregator.model.Parameter;
import gios.data.model.MeasurementDTO;
import openaq.data.model.Measurement;

import java.util.List;

/**
 * Helper class to hold parameter info with associated measurements
 */
public record ParameterWithMeasurements(Parameter parameter, List<MeasurementDTO> giosMeasurements,
                                        List<Measurement> openaqMeasurements) {
}
