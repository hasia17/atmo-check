package gios_data.domain.repository.station;

import gios_data.domain.model.Station;

import java.util.List;

public interface StationRepositoryCustom {
    List<Station> searchByCriteria(String city, String province, Double lat, Double lon, Double radiusKm);
}