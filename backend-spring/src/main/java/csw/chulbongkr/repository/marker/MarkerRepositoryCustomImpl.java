package csw.chulbongkr.repository.marker;

import csw.chulbongkr.dto.MarkerDTO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

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

    @Override
    public List<MarkerDTO.MarkerWithDistance> findMarkersWithinDistance(double latitude, double longitude, double distance, int pageSize, int offset) {
        String point = String.format("POINT(%f %f)", latitude, longitude);
        String query = """
            SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude, Description,\s
                   ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) AS distance, Address
            FROM Markers
            WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?
            ORDER BY distance
            LIMIT ? OFFSET ?
           \s""";

        return jdbcTemplate.query(query,
                (rs, rowNum) -> new MarkerDTO.MarkerWithDistance(
                        rs.getInt("MarkerID"),
                        rs.getDouble("Latitude"),
                        rs.getDouble("Longitude"),
                        rs.getString("Description"),
                        rs.getDouble("distance"),
                        Optional.ofNullable(rs.getString("Address"))
                ),
                point, point, distance, pageSize, offset);
    }
}
