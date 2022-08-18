package io.pirs.data.datamgmt.config;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.pirs.data.datamgmt.proto.TrackerGrpc;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GrpcConfig {

    @Bean
    public TrackerGrpc.TrackerBlockingStub blockingTrackerStub() {
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 8080)
                .usePlaintext()
                .build();

        return TrackerGrpc.newBlockingStub(channel);
    }
}
