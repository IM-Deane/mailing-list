# Project: Golang Mailing list CLI

### Overview

This project is a Go CLI microservice that allows a user to interact with a
mailing list database.

Requests can be made via two different API endpoints: HTTP or gRPC.

### Tech Stack:

- SQLite for data storage
- Protocol Buffers for communication format
- gRPC client needed to interact with gRPC server

### Running the project

You can start the gRPC and JSON servers with: `go run ./server`

You can start the gRPC client with: `go run ./client`

### Testing the project:

**JSON:** You can test the JSON API with cURL, Postman, or Thunder Client (a VS
Code extension)

Example endpoint:

`http://127.0.0.1:8080/email/get`

Can fetch an email from the database using this endpoint.

**gRPC:** You can test the gRPC server via the gRPC client.

`./client/client.go, line 124` has several test requests

## Development Setup

If you wish to fork or edit this project, it requires a `gcc` compiler installed
and the `protobuf` code generation tools.

### Install protobuf compiler

Install the `protoc` tool using the instructions available at
[https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/).

Alternatively you can download a pre-built binary from
[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
and placing the extracted binary somewhere in your `$PATH`.

### Install Go protobuf codegen tools

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

### Generate Go code from .proto files

```
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  Proto/mail.proto
```
