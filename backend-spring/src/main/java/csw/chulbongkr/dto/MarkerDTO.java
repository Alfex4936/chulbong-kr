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
            int markerId,
            double latitude,
            double longitude,
            String description
            // String username,
            // int userId,
            // List<String> photoUrls
    ) {}

    public record MarkerWithDistance(
            int markerId,
            double latitude,
            double longitude,
            String description,
            double distance,
            Optional<String> address
    ) {}

    public record MarkerWithDislike(
            int markerId,
            double latitude,
            double longitude,
            String username,
            int dislikeCount
    ) {}

    public record MarkerSimple(
            int markerId,
            double latitude,
            double longitude
    ) {}

    public record MarkerSimpleWithDescription(
            int markerId,
            double latitude,
            double longitude,
            String description,
            LocalDateTime createdAt,
            Optional<String> address
    ) {}

    public record MarkerSimpleWithAddr(
            int markerId,
            double latitude,
            double longitude,
            Optional<String> address
    ) {}
}
