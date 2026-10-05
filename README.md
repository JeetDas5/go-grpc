# go-grpc project Go

## Overview

This repository contains a simple **gRPC** example written in Go. The project is structured into three main components:

- **`proto/`** – Protocol Buffer definitions for the service and messages.
- **`server/`** – Implementation of the gRPC server exposing the `GreetService`.
- **`client/`** – A client that demonstrates unary, server‑streaming, client‑streaming, and bidirectional‑streaming calls.

The design follows a classic client‑server architecture where the server provides the service and the client consumes it. All code lives under a single Go module (defined in `go.mod`).

## Directory Layout

```
├── client/                # Client application source
│   └── main.go
├── server/                # Server application source
│   └── main.go
├── proto/                 # .proto files (gRPC contract)
│   ├── greet.proto
│   ├── greet.pb.go
│   └── greet_grpc.pb.go
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
└── README.md              # This documentation
```

## Prerequisites

- Go 1.22 or later installed ([download](https://go.dev/dl/)).
- `protoc` (Protocol Buffers compiler) installed.
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

These plugins must be available in your `PATH` for code generation.

## Setup & Build

1. **Clone the repository** (if you haven't already):

```bash
git clone https://github.com/JeetDas5/go-grpc.git
cd go-grpc
```

2. **Generate Go code from the protobuf definitions**:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/greet.proto
```

   This command creates `proto/greet.pb.go` and `proto/greet_grpc.pb.go`.

3. **Download Go module dependencies**:

```bash
go mod tidy
```

## Running the Server

From the project root, start the gRPC server:

```bash
go run ./server
```

The server listens on **`localhost:8080`**.

## Running the Client

In a separate terminal, execute the client to make a bidirectional‑streaming call (default behaviour in `client/main.go`):

```bash
go run ./client
```

You can modify `client/main.go` to call other RPC methods (unary, server‑streaming, client‑streaming) by uncommenting the corresponding function calls.

## Testing the Service

The repository includes a basic test suite. Run the tests with:

```bash
go test ./...
```

Built by: [Jeet Das](https://github.com/JeetDas5)
