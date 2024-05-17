package csw.chulbongkr.repository.marker;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.entity.Marker;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface MarkerRepository extends JpaRepository<Marker, Integer>, MarkerRepositoryCustom {
    List<MarkerDTO.MarkerSimple> findAllSimplifiedMarkers();
    List<MarkerDTO.MarkerWithDistance> findMarkersWithinDistance(double latitude, double longitude, double distance, int pageSize, int offset);
}
