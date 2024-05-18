package csw.chulbongkr.dto;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.time.LocalDateTime;
import java.util.Optional;

public class MarkerDTO {

    @JsonInclude(JsonInclude.Include.NON_ABSENT) // Omit fields that are absent (Optional.empty())
    public record MarkerRequest(
            Optional<Integer> markerId,
            double latitude,
            double longitude,
            String description,
            Optional<String> photoUrl
    ) {}

    public record MarkerResponse(
            Integer markerId,
            double latitude,
            double longitude,
            String description
            // String username,
            // Integer userId,
            // List<String> photoUrls
    ) {}

    @JsonInclude(JsonInclude.Include.NON_ABSENT)
    public record MarkerWithDistance(
            Integer markerId,
            double latitude,
            double longitude,
            String description,
            double distance,
            Optional<String> address
    ) {}

    public record MarkerWithDislike(
            Integer markerId,
            double latitude,
            double longitude,
            String username,
            Integer dislikeCount
    ) {}

    public record MarkerSimple(
            Integer markerId,
            double latitude,
            double longitude
    ) {}

    public record MarkerSimpleWithDescription(
            Integer markerId,
            double latitude,
            double longitude,
            String description,
            LocalDateTime createdAt,
            Optional<String> address
    ) {}

    public record MarkerSimpleWithAddr(
            Integer markerId,
            double latitude,
            double longitude,
            Optional<String> address
    ) {}

    public record FindMarkerNearbyQuery(
            double latitude,
            double longitude,
            double distance,
            int pageSize,
            int offset
    ) {}
}
