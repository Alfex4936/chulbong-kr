package csw.chulbongkr.config.custom;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Getter
@Setter
@Configuration
@ConfigurationProperties(prefix = "slack")
public class SlackConfig {
    private String botToken;
    private String channelId;
}