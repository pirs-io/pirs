package com.pirs.gateway.config.route;

import org.apache.camel.builder.RouteBuilder;
import org.springframework.stereotype.Component;

@Component
public class GrpcRoutesConfig extends RouteBuilder {

    @Override
    public void configure() throws Exception {

        from("restApi:")
                .to("grpc:");

    }
}
