package csw.chulbongkr.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

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

    public record Weather(
            Codes codes,
            WeatherInfos weatherInfos
    ) {
        public record Codes(
                @JsonProperty("resultCode") String resultCode,
                @JsonProperty("hcode") Code hcode,
                @JsonProperty("bcode") Code bcode
        ) {}

        public record Code(
                @JsonProperty("type") String type,
                @JsonProperty("code") String code,
                @JsonProperty("name") String name,
                @JsonProperty("fullName") String fullName,
                @JsonProperty("regionId") String regionId,
                @JsonProperty("name0") String name0,
                @JsonProperty("code1") String code1,
                @JsonProperty("name1") String name1,
                @JsonProperty("code2") String code2,
                @JsonProperty("name2") String name2,
                @JsonProperty("code3") String code3,
                @JsonProperty("name3") String name3,
                @JsonProperty("childcount") float childcount,
                @JsonProperty("x") float x,
                @JsonProperty("y") float y
        ) {}

        public record WeatherInfos(
                @JsonProperty("current") WeatherInfo current,
                @JsonProperty("forecast") WeatherInfo forecast
        ) {}

        public record WeatherInfo(
                @JsonProperty("type") String type,
                @JsonProperty("rcode") String rcode,
                @JsonProperty("iconId") String iconId,
                @JsonProperty("temperature") String temperature,
                @JsonProperty("desc") String desc,
                @JsonProperty("humidity") String humidity,
                @JsonProperty("rainfall") String rainfall,
                @JsonProperty("snowfall") String snowfall
        ) {}

        public record WeatherRequest(
                @JsonProperty("temperature") String temperature,
                @JsonProperty("desc") String desc,
                @JsonProperty("iconImage") String iconImage,
                @JsonProperty("humidity") String humidity,
                @JsonProperty("rainfall") String rainfall,
                @JsonProperty("snowfall") String snowfall
        ) {}
    }
}
