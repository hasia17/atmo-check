package aggregator.service.wrappers;

import java.time.OffsetDateTime;

/**
 * Helper class to hold aggregated measurement statistics
 */
public record AggregatedMeasurements(Double averageValue, Double minValue, Double maxValue, Integer measurementCount,
                                     Double latestValue, OffsetDateTime latestTimestamp) {
}
