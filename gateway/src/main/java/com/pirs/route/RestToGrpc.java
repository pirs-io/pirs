package com.pirs.route;

import org.apache.camel.Exchange;
import org.apache.camel.Expression;
import org.apache.camel.builder.RouteBuilder;

import javax.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class RestToGrpc extends RouteBuilder {

    @Override
    public void configure() throws Exception {
        from("direct:getUserByName")
                .transform(new Expression() {
                    @Override
                    public <T> T evaluate(Exchange exchange, Class<T> type) {
                        return null;
                    }
                }).end();
    }
}
