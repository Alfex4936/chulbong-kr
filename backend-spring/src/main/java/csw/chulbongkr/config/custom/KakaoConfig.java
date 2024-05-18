package csw.chulbongkr.config.custom;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Getter
@Setter
@Configuration
@ConfigurationProperties(prefix = "kakao")
public class KakaoConfig {
    private String staticMap;
    private String restApiKey;

    private String waterApiKey;
    private String waterApiUrl;

    private String weatherUrl;
    private String weatherIconUrl;

    private String addressInfo;

    private String coordAddr;
    private String coordRegion;
    private String geoCode;
}