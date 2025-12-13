package open.meteo.domain.repository;

import open.meteo.domain.model.Parameter;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ParameterRepository extends MongoRepository<Parameter, String> {
}
