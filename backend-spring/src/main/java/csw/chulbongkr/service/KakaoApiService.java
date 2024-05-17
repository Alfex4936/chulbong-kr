package csw.chulbongkr.service;

import csw.chulbongkr.config.custom.KakaoConfig;
import csw.chulbongkr.dto.KakaoDTO;
import csw.chulbongkr.util.CoordinatesConverter;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.Optional;

@RequiredArgsConstructor
@Service
public class KakaoApiService {
    private final KakaoConfig kakaoConfig;

    private final RestTemplate restTemplate;

    // kakao map way to fetch address
    public String fetchAddress(double lat, double lng) {
        CoordinatesConverter.XYCoordinate xy = CoordinatesConverter.convertWGS84ToWCONGNAMUL(lat, lng);
        String requestURL = kakaoConfig.getAddressInfo() + "&x=" + xy.latitude() + "&y=" + xy.longitude();

        KakaoDTO.KakaoMarkerData apiResp = restTemplate.getForObject(requestURL, KakaoDTO.KakaoMarkerData.class);

        return Optional.ofNullable(apiResp)
                .flatMap(this::getAddress)
                .orElse("대한민국 철봉 지도");
    }

    private Optional<String> getAddress(KakaoDTO.KakaoMarkerData apiResp) {
        return getNewAddress(apiResp).or(() -> getOldAddress(apiResp));
    }

    private Optional<String> getNewAddress(KakaoDTO.KakaoMarkerData apiResp) {
        return Optional.ofNullable(apiResp.newAddr())
                .flatMap(newAddr -> {
                    if (newAddr.name() != null && !newAddr.name().isEmpty()) {
                        StringBuilder address = new StringBuilder(newAddr.name());
                        if (newAddr.building() != null && !newAddr.building().isEmpty()) {
                            address.append(", ").append(newAddr.building());
                        }
                        return Optional.of(address.toString());
                    }
                    return Optional.empty();
                });
    }

    private Optional<String> getOldAddress(KakaoDTO.KakaoMarkerData apiResp) {
        return Optional.ofNullable(apiResp.old())
                .flatMap(oldAddr -> {
                    if (oldAddr.name() != null && !oldAddr.name().isEmpty()) {
                        StringBuilder address = new StringBuilder(oldAddr.name());
                        if (oldAddr.building() != null && !oldAddr.building().isEmpty()) {
                            address.append(", ").append(oldAddr.building());
                        }
                        return Optional.of(address.toString());
                    }
                    return Optional.empty();
                });
    }

    public String getStaticImageURL() {
        return kakaoConfig.getStaticMap();
    }
}
