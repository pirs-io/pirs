package io.pirs.data.datamgmt;

import io.pirs.data.datamgmt.proto.LocationRequest;
import io.pirs.data.datamgmt.proto.PackageLocation;
import io.pirs.data.datamgmt.proto.TrackerGrpc;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RestController
public class DataMgmtApplication {

    @Autowired
    private TrackerGrpc.TrackerBlockingStub trackerBlockingStub;

    public static void main(String[] args) {
        SpringApplication.run(DataMgmtApplication.class, args);
    }

    @GetMapping
    public String run() throws Exception {
        PackageLocation packageLocation = trackerBlockingStub.findPackageLocation(LocationRequest.newBuilder().build());
        return packageLocation.toString();
    }
}
