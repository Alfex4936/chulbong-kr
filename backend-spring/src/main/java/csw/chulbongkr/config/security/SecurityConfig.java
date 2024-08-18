package csw.chulbongkr.config.security;

import csw.chulbongkr.service.auth.CustomUserDetailsService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.dao.DaoAuthenticationProvider;
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.annotation.web.configurers.HeadersConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.core.session.SessionRegistry;
import org.springframework.security.core.session.SessionRegistryImpl;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.security.web.firewall.StrictHttpFirewall;
import org.springframework.security.web.header.Header;
import org.springframework.security.web.header.writers.*;
import org.springframework.web.cors.CorsConfigurationSource;

import java.util.List;

import static org.springframework.security.web.util.matcher.AntPathRequestMatcher.antMatcher;

@Slf4j
@Configuration
@EnableMethodSecurity
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfig {

    private final CustomUserDetailsService userDetailService;

    @Qualifier("corsConfigurationSource")
    private final CorsConfigurationSource corsConfigurationSource;

    private final JwtTokenProvider jwtTokenProvider;

    @Bean
    public StrictHttpFirewall httpFirewall() {
        StrictHttpFirewall firewall = new StrictHttpFirewall();
        firewall.setAllowSemicolon(true);
        firewall.setAllowBackSlash(true);
        firewall.setAllowUrlEncodedDoubleSlash(true);
        return firewall;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        // Disable CSRF as API is stateless
        http.csrf(AbstractHttpConfigurer::disable);
        http.anonymous(AbstractHttpConfigurer::disable);

        // Configure CORS using a custom configuration source
        http.cors(corsConfigurer -> corsConfigurer.configurationSource(corsConfigurationSource));

        // Disable session creation as API is stateless
        http.sessionManagement(session -> session
                .sessionCreationPolicy(SessionCreationPolicy.STATELESS)
                .maximumSessions(1)
                .sessionRegistry(sessionRegistry()));

        // @formatter:off
        http.headers(
                headers ->
                        headers
                                .cacheControl(HeadersConfigurer.CacheControlConfig::disable)
                                .addHeaderWriter(new StaticHeadersWriter(
                                        List.of(
                                                new Header("X-Dns-Prefetch-Control", "on"),
                                                new Header("X-Download-Options", "noopen"),
                                                new Header("X-Permitted-Cross-Domain-Policies", "none")
                                        )))
                                .crossOriginEmbedderPolicy(coep -> coep.policy(CrossOriginEmbedderPolicyHeaderWriter.CrossOriginEmbedderPolicy.REQUIRE_CORP))
                                .crossOriginOpenerPolicy(coop -> coop.policy(CrossOriginOpenerPolicyHeaderWriter.CrossOriginOpenerPolicy.SAME_ORIGIN))
                                .crossOriginResourcePolicy(corp -> corp.policy(CrossOriginResourcePolicyHeaderWriter.CrossOriginResourcePolicy.SAME_ORIGIN))
                                .frameOptions(HeadersConfigurer.FrameOptionsConfig::sameOrigin)
                                .referrerPolicy(referrerPolicy -> referrerPolicy.policy(ReferrerPolicyHeaderWriter.ReferrerPolicy.SAME_ORIGIN))
                                .httpStrictTransportSecurity(hstsConfig -> hstsConfig.maxAgeInSeconds(31536000).includeSubDomains(true))
                                .xssProtection(xss -> xss.headerValue(XXssProtectionHeaderWriter.HeaderValue.ENABLED_MODE_BLOCK))
                                .contentSecurityPolicy(cps -> cps.policyDirectives("default-src 'self';base-uri 'self';" + "font-src 'self' https: data:;form-action 'self';" + "frame-ancestors 'self';img-src 'self' data:;" + "object-src 'none';script-src 'self';" + "script-src-attr 'none';style-src 'self' https: 'unsafe-inline';" + "upgrade-insecure-requests")
                                ));

        http.authorizeHttpRequests(auth -> auth.requestMatchers(HttpMethod.OPTIONS, "/*").permitAll() // Allow preflight requests for all paths
                .requestMatchers(antMatcher("/")).permitAll() // Allow all requests to the root path
                .requestMatchers(antMatcher("/api/v1/users/**")).authenticated() // Secure all API endpoints
                .requestMatchers(antMatcher("/api/v1/**")).permitAll() // Secure all API endpoints
                .anyRequest().permitAll() // Allow all other requests
        );

        http.addFilterBefore(authenticationFilter(), UsernamePasswordAuthenticationFilter.class);
        return http.build();
    }

    @Bean
    public AuthenticationFilter authenticationFilter() {
        return new AuthenticationFilter(jwtTokenProvider, userDetailService);
    }

    @Bean
    public UserDetailsService userDetailsService() {
        return userDetailService;
    }

    @Bean
    public BCryptPasswordEncoder bCryptPasswordEncoder() {
        return new BCryptPasswordEncoder();
    }


    @Bean
    public SessionRegistry sessionRegistry() {
        return new SessionRegistryImpl();
    }

    @Bean
    public AuthenticationManager authenticationManager(AuthenticationConfiguration authConfig) throws Exception {
        return authConfig.getAuthenticationManager();
    }

    @Bean
    public DaoAuthenticationProvider authenticationProvider() {
        DaoAuthenticationProvider authProvider = new DaoAuthenticationProvider();

        authProvider.setUserDetailsService(userDetailService);
        authProvider.setPasswordEncoder(bCryptPasswordEncoder());

        return authProvider;
    }
}
