package csw.chulbongkr.config;

import com.alibaba.fastjson2.JSONReader;
import com.alibaba.fastjson2.JSONWriter;
import com.alibaba.fastjson2.filter.SimplePropertyPreFilter;
import com.alibaba.fastjson2.support.config.FastJsonConfig;
import com.alibaba.fastjson2.support.spring6.http.converter.FastJsonHttpMessageConverter;
import org.springframework.boot.autoconfigure.http.HttpMessageConverters;
import org.springframework.context.annotation.Bean;
import org.springframework.http.MediaType;
import org.springframework.http.converter.HttpMessageConverter;

import java.nio.charset.StandardCharsets;
import java.util.Collections;

public class CustomFastJsonConfig {

    @Bean
    public FastJsonHttpMessageConverter fastJsonHttpMessageConverter() {
        FastJsonHttpMessageConverter converter = new FastJsonHttpMessageConverter();

        FastJsonConfig config = new FastJsonConfig();
        config.setDateFormat("yyyy-MM-dd'T'HH:mm:ss");
        config.setReaderFeatures(JSONReader.Feature.FieldBased, JSONReader.Feature.SupportArrayToBean);
        config.setCharset(StandardCharsets.UTF_8);
        config.setWriterFeatures(JSONWriter.Feature.WriteMapNullValue);

        // Filter to exclude null values and handle Optional fields
        SimplePropertyPreFilter filter = new SimplePropertyPreFilter();
        filter.getExcludes().add("null");
        config.setWriterFilters(filter);

        converter.setDefaultCharset(StandardCharsets.UTF_8);
        converter.setFastJsonConfig(config);
        converter.setSupportedMediaTypes(Collections.singletonList(MediaType.APPLICATION_JSON));
        return converter;
    }

    @Bean
    public HttpMessageConverters customConverters() {
        HttpMessageConverter<?> additional = fastJsonHttpMessageConverter();
        return new HttpMessageConverters(additional);
    }
}
