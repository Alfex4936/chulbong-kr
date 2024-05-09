package csw.chulbongkr.service;

import csw.chulbongkr.dto.MarkerDTO;
import csw.chulbongkr.repository.marker.MarkerRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@RequiredArgsConstructor
@Service
public class MarkerService {
    private final MarkerRepository markerRepository;

    public List<MarkerDTO.MarkerSimple> getAllMarkers() {
        return markerRepository.findAllSimplifiedMarkers();
    }

}
