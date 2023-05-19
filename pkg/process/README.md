# Process
A microservice for process file management.

## How to run

### Requirements
In order to fully run this microservice, you have to fulfill some **requirements**:
1. Running instance of *MongoDB* database
2. Running instance of *Dependency-management* microservice. See [README](../dependency-management/README.md).
3. Running instance of *Process-storage* microservice. See [README](../process-storage/README.md).

### Configuration
When you have the requirements fulfilled, you have to set up the configuration for this microservice. Application
expects *'.env'* file on root folder. The other option for configuration is to set up environment variables. They have
precedence over *'.env'* file. An example config file can be found [here](example.env).

### Startup
You can run the microservice by running `main.go` file or you can use the [Dockerfile](Dockerfile) to create container and 
run it without any args.

## How to use
This microservice publishes API. To use this API we have created temporary client tool for development use. It can be
found [here](utils/process_client.go).

You must specify the `host` and `port` of `Process-service` in constants. Then you can run the file and use command-line
menu to navigate.