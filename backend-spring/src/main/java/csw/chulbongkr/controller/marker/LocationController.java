package csw.chulbongkr.controller.marker;

import csw.chulbongkr.service.KakaoApiService;
import csw.chulbongkr.service.MarkerService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/markers")
public class LocationController {
    private final MarkerService markerService;
    private final KakaoApiService kakaoApiService;

    @GetMapping("/test")
    public String test() {
        return "test";
    }
}
