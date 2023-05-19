# Dependency-management-service
A microservice for process dependency management.

## How to run

### Requirements
In order to fully run this microservice, you have to have running instance of *MongoDB* database.

### Configuration
When you have the requirements fulfilled, you have to set up the configuration for this microservice. Application
expects *'.env'* file on root folder. The other option for configuration is to set up environment variables. They have
precedence over *'.env'* file. An example config file can be found [here](example.env).

### Startup
You can run the microservice by running `main.go` file.

## How to use
This microservice is used by *Process-service* microservice automatically.
