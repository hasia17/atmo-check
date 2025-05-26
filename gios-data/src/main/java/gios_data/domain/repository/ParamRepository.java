package gios_data.domain.repository;

import gios_data.domain.model.Param;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ParamRepository extends MongoRepository<Param, Integer> {
}
