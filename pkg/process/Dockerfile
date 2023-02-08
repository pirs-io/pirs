FROM golang:1.19-bullseye as builder

ARG APP_NAME=process
ARG APP_DIR=/app/$APP_NAME
ARG COMMONS_DIR=/app/commons
ARG PROCESS_STORAGE_DIR=/app/process-storage

WORKDIR /app

## install protoc
RUN apt-get update
RUN apt-get install -y protobuf-compiler

# install go protoc deps
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/envoyproxy/protoc-gen-validate@latest

# copy all application files and dependencies
COPY $APP_NAME $APP_DIR
COPY commons $COMMONS_DIR
COPY process-storage $PROCESS_STORAGE_DIR

COPY $APP_NAME/go.work-template .
RUN mv go.work-template go.work

# copy all needed .proto files from other services
COPY ./process-storage/grpc $APP_DIR/grpc
# generate go files from proto
RUN cd $APP_DIR/grpc && ./generate.sh

WORKDIR $APP_DIR

# install deps && build app
RUN go work sync
RUN go work use .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$APP_NAME .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/process/bin/process .
COPY --from=builder /app/process/example.env .

CMD ["./process"]