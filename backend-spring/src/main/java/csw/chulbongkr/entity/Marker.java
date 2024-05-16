package csw.chulbongkr.entity;

import com.fasterxml.jackson.annotation.JsonFormat;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import org.locationtech.jts.geom.Point;

import java.time.LocalDateTime;

@Entity
@Table(name = "Markers")
public class Marker {

    @Id
    @Column(name = "MarkerID")
    private Integer markerID;

    @Column(name = "UserID")
    private Integer userID;

    @Column(columnDefinition = "geometry(Point, 4326)")
    private Point location;

//    @Column(name = "Latitude")
//    private Double latitude;
//
//    @Column(name = "Longitude")
//    private Double longitude;

    @Column(name = "Description")
    private String description;

    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    @Column(name = "CreatedAt")
    private LocalDateTime createdAt;

    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    @Column(name = "UpdatedAt")
    private LocalDateTime updatedAt;

    @Column(name = "Address")
    private String address;

}
