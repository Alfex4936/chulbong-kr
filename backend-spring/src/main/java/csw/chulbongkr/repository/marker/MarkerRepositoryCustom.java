package csw.chulbongkr.repository.marker;

import csw.chulbongkr.dto.MarkerDTO;

import java.util.List;

public interface MarkerRepositoryCustom {
    List<MarkerDTO.MarkerSimple> findAllSimplifiedMarkers();
    List<MarkerDTO.MarkerSimpleWithAddr> findAllSimplifiedMarkersWithAddress();
    List<MarkerDTO.MarkerWithDistance> findMarkersWithinDistance(double latitude, double longitude, double distance, int pageSize, int offset);
}
