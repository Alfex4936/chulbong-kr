package csw.chulbongkr.dto;

import java.util.List;

public class KakaoDTO {
    public record Coordinate(double x, double y) {}

    public record Shape(boolean hole, String type, List<List<Coordinate>> coordinateList) {}

    public record AreaDocument(
            String docid,
            String name,
            String roadName,
            String building,
            String bunji,
            String ho,
            String san,
            Shape shape,
            String zoneNo) {}

    public record KakaoMarkerData(
            AreaDocument old,
            double x,
            double y,
            String regionid,
            String region,
            String bcode,
            String linePtsFormat,
            AreaDocument newAddr) {}
}
