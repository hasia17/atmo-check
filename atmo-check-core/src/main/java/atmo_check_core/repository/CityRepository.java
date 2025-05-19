package atmo_check_core.repository;

import atmo_check_core.model.City;
import atmo_check_core.model.Param;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface CityRepository extends MongoRepository<City, Integer> {
}
