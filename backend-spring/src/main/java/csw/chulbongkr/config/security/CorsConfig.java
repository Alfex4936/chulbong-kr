package csw.chulbongkr.config.security;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.http.HttpMethod;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import org.springframework.web.filter.CorsFilter;

import java.util.Arrays;
import java.util.List;

@Configuration
public class CorsConfig {

    @Bean
    public CorsFilter corsFilter() {
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", corsConfiguration());

        return new CorsFilter(source);
    }

    private CorsConfiguration corsConfiguration() {
        CorsConfiguration configuration = new CorsConfiguration();
        configuration.setAllowCredentials(true);
        configuration.setAllowedMethods(Arrays.asList(HttpMethod.GET.name(), HttpMethod.POST.name(), HttpMethod.PUT.name(), HttpMethod.DELETE.name(), HttpMethod.OPTIONS.name()));
        configuration.setAllowedHeaders(List.of("*"));
        configuration.setExposedHeaders(List.of("Accept"));

        // Add specific non-wildcard origins
        configuration.setAllowedOrigins(Arrays.asList("https://chulbong-kr.vercel.app", "http://localhost:5173"));

        return configuration;
    }

    @Bean
    @Primary
    public CorsConfigurationSource corsConfigurationSource() {
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", corsConfiguration());
        return request -> {
            String origin = request.getHeader("Origin");
            if (origin != null) {
                if (origin.endsWith(".k-pullup.com")) {
                    CorsConfiguration config = new CorsConfiguration(corsConfiguration());
                    config.setAllowedOrigins(List.of(origin));
                    return config;
                }
            }
            return corsConfiguration();
        };
    }
}
