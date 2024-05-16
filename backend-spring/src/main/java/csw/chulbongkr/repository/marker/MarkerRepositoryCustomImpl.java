package csw.chulbongkr.repository.marker;

import csw.chulbongkr.dto.MarkerDTO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public class MarkerRepositoryCustomImpl implements MarkerRepositoryCustom {

    private final JdbcTemplate jdbcTemplate;

    @Autowired
    public MarkerRepositoryCustomImpl(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    @Override
    public List<MarkerDTO.MarkerSimple> findAllSimplifiedMarkers() {
        String sql = "SELECT MarkerID, ST_X(Location) AS latitude, ST_Y(Location) AS longitude FROM Markers";
        return jdbcTemplate.query(
                "SELECT MarkerID, ST_X(Location) as Latitude, ST_Y(Location) as Longitude FROM Markers",
                (rs, rowNum) -> new MarkerDTO.MarkerSimple(
                        rs.getInt("MarkerID"),
                        rs.getDouble("Latitude"),
                        rs.getDouble("Longitude")
                )
        );
    }
}
