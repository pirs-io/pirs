package com.pirs.gateway.config.rest;

import org.apache.camel.CamelContext;
import org.apache.camel.Component;
import org.apache.camel.component.rest.openapi.RestOpenApiComponent;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.net.URI;
import java.net.URISyntaxException;

@Configuration
public class RestApi {

    @Bean
    public Component restApi(CamelContext camelContext) throws URISyntaxException {
        RestOpenApiComponent petstore = new RestOpenApiComponent(camelContext);
        petstore.setSpecificationUri(new URI("classpath:openapi.json"));
        petstore.setHost("0.0.0.0");
        return petstore;
    }
}
