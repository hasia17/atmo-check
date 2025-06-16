package gios_data.domain.repository.station;

import gios_data.domain.model.Station;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.query.Criteria;

import org.springframework.data.mongodb.core.query.Query;
import org.springframework.stereotype.Repository;

import java.util.List;

@Slf4j
@Repository
public class StationRepositoryImpl implements StationRepositoryCustom {

    private final MongoTemplate mongoTemplate;

    public StationRepositoryImpl(MongoTemplate mongoTemplate) {
        this.mongoTemplate = mongoTemplate;
    }

    @Override
    public List<Station> searchByCriteria(String city, String province, Double lat, Double lon, Double radiusKm) {
        Query query = new Query();

        if (!StringUtils.isEmpty(city)) {
            query.addCriteria(Criteria.where("city").regex("^" + city + "$", "i"));
        }
        if (!StringUtils.isEmpty(province)) {
            query.addCriteria(Criteria.where("province").regex("^" + province + "$", "i"));
        }
        if (lat != null && lon != null && radiusKm != null) {
            if (lat < -90 || lat > 90) {
                throw new IllegalArgumentException("Latitude must be between -90 and 90");
            }
            if (lon < -180 || lon > 180) {
                throw new IllegalArgumentException("Longitude must be between -180 and 180");
            }
            double latRadius = radiusKm / 111.0; // 1 deg latitude ≈ 111 km
            double lonRadius = radiusKm / (111.0 * Math.cos(Math.toRadians(lat)));
            query.addCriteria(Criteria.where("gegrLat").gte(lat - latRadius).lte(lat + latRadius));
            query.addCriteria(Criteria.where("gegrLon").gte(lon - lonRadius).lte(lon + lonRadius));
        } else {
            log.info("All three: lat, lon and radius must be provided to search by geographical coordinates");
        }

        return mongoTemplate.find(query, Station.class);
    }
}