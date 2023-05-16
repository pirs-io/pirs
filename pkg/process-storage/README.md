# storage-service

service providing storage for process files

currently supported storage providers:
    * git

## running this service


**configuration**
application reads .env file on root of this project containing all configuraiton properties
all properties that are specified as env variables have precedence over .env file
all properties in .exampleenv are required for application to run

**running**
run main.go or build Dockerfile and run image without any args


## code organization

when contributing follow this code structure, if any feature does not fit in any package, create new package

### adapter
contains definition of storage adapter interface along with implementation of supported concrete adapters
every new adapter implementation must be placed here with _adapter filename suffix

### config
contains all configuration logic needed for application to run
initialize all singleton service instances here and put them into ```ApplicationContext```

### grpc
place for .proto file and generated grpc code

### storage
implementation of storage providers


### root 
```main.go``` - running application
```grpc_server.go``` - grpc server implementation

