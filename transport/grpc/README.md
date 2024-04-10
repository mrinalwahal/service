# gRPC Server

## Prerequisites

Install gRPC protoc tools from [here](https://grpc.io/docs/languages/go/quickstart/#prerequisites).

## Usage

Generate protobuf by running the following commands in this directory.

```bash
export PATH="${PATH}:${HOME}/go/bin"
```

```bash
protoc --go_out=gen --go_opt=paths=source_relative --go-grpc_out=gen --go-grpc_opt=paths=source_relative service.proto
```