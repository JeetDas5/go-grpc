# go-grpc

A reference implementation of all four gRPC communication patterns in Go, built around a `GreetService` defined with Protocol Buffers.

---

## Table of Contents

- [Overview](#overview)
- [Technical Architecture](#technical-architecture)
- [Repository Layout](#repository-layout)
- [Service Contract](#service-contract)
- [RPC Patterns](#rpc-patterns)
  - [Unary](#1-unary)
  - [Server Streaming](#2-server-streaming)
  - [Client Streaming](#3-client-streaming)
  - [Bidirectional Streaming](#4-bidirectional-streaming)
- [Prerequisites](#prerequisites)
- [Code Generation](#code-generation)
- [Running](#running)
- [Dependencies](#dependencies)

---

## Overview

This project demonstrates how gRPC works in Go across all streaming variants. The server implements a single `GreetService` and registers all four RPCs. The client connects over an insecure channel and exercises each call pattern independently through separate, self-contained functions.

---

## Technical Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Client (client/)                  │
│                                                      │
│  main.go          – dials localhost:8080, builds    │
│                     pb.GreetServiceClient           │
│  unary.go         – callSayHello()                  │
│  server_stream.go – callSayHelloServerStream()      │
│  client_stream.go – callSayHelloClientStream()      │
│  bi_stream.go     – callSayHelloBiDirectionalStreaming() │
└──────────────────────┬──────────────────────────────┘
                       │  HTTP/2  (insecure, port 8080)
                       │  Protocol Buffers (wire format)
┌──────────────────────▼──────────────────────────────┐
│                    Server (server/)                  │
│                                                      │
│  main.go          – net.Listen(tcp, :8080),         │
│                     grpc.NewServer(),               │
│                     RegisterGreetServiceServer()    │
│  unary.go         – SayHello()                      │
│  server_stream.go – SayHelloServerStreaming()       │
│  client_stream.go – SayHelloClientStreaming()       │
│  bi_stream.go     – SayHelloBiDirectionalStreaming()│
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│               Proto Layer (proto/)                   │
│                                                      │
│  greet.proto       – source of truth (IDL)          │
│  greet.pb.go       – generated message types        │
│  greet_grpc.pb.go  – generated client & server stubs│
└─────────────────────────────────────────────────────┘
```

The proto layer is the single source of truth. Both the server and the client import only the generated `proto` package; neither depends directly on the other. The transport is raw TCP with HTTP/2 multiplexing handled by `google.golang.org/grpc`. No TLS is configured — the client dials with `insecure.NewCredentials()`.

---

## Repository Layout

```
go-grpc/
├── client/
│   ├── main.go              # Entry point: dials server, selects RPC call
│   ├── unary.go             # Unary call implementation
│   ├── server_stream.go     # Server-streaming call implementation
│   ├── client_stream.go     # Client-streaming call implementation
│   └── bi_stream.go         # Bidirectional-streaming call implementation
├── server/
│   ├── main.go              # Entry point: TCP listener, gRPC server setup
│   ├── unary.go             # SayHello handler
│   ├── server_stream.go     # SayHelloServerStreaming handler
│   ├── client_stream.go     # SayHelloClientStreaming handler
│   └── bi_stream.go         # SayHelloBiDirectionalStreaming handler
├── proto/
│   ├── greet.proto          # Protobuf IDL (edit this, regenerate the rest)
│   ├── greet.pb.go          # Auto-generated: message structs
│   └── greet_grpc.pb.go     # Auto-generated: client stub + server interface
├── go.mod
├── go.sum
└── README.md
```

---

## Service Contract

Defined in [`proto/greet.proto`](proto/greet.proto):

```protobuf
service GreetService {
  rpc SayHello(NoParam)                                         returns (HelloResponse);
  rpc SayHelloServerStreaming(NamesList)                        returns (stream HelloResponse);
  rpc SayHelloClientStreaming(stream HelloRequest)              returns (MessagesList);
  rpc SayHelloBiDirectionalStreaming(stream HelloRequest)       returns (stream HelloResponse);
}

message NoParam       {}
message HelloRequest  { string name     = 1; }
message HelloResponse { string message  = 1; }
message NamesList     { repeated string names    = 1; }
message MessagesList  { repeated string messages = 1; }
```

---

## RPC Patterns

### 1. Unary

```
Client                          Server
  |-- SayHello(NoParam) ------->|
  |<-- HelloResponse -----------|
```

- Client sends a single request with no parameters.
- Server returns a single `HelloResponse{Message: "Hello World!"}`.
- Client uses a `context.WithTimeout` of 1 second.

**Files:** [`client/unary.go`](client/unary.go), [`server/unary.go`](server/unary.go)

---

### 2. Server Streaming

```
Client                             Server
  |-- SayHelloServerStreaming(NamesList) -->|
  |<-- HelloResponse (name[0]) ------------|
  |<-- HelloResponse (name[1]) ------------|
  |<-- HelloResponse (name[n]) ------------|
  |<-- EOF --------------------------------|
```

- Client sends a `NamesList` once.
- Server iterates the names, sends one `HelloResponse` per name with a 2-second delay between each, then closes the stream.
- Client reads with `stream.Recv()` in a loop until `io.EOF`.

**Files:** [`client/server_stream.go`](client/server_stream.go), [`server/server_stream.go`](server/server_stream.go)

---

### 3. Client Streaming

```
Client                              Server
  |-- HelloRequest (name[0]) ------>|
  |-- HelloRequest (name[1]) ------>|
  |-- HelloRequest (name[n]) ------>|
  |-- CloseAndRecv() -------------->|
  |<-- MessagesList ----------------|
```

- Client opens a stream, sends one `HelloRequest` per name (2-second delay between sends), then calls `CloseAndRecv()`.
- Server accumulates each greeting into a `[]string`, then returns a single `MessagesList` via `SendAndClose`.

**Files:** [`client/client_stream.go`](client/client_stream.go), [`server/client_stream.go`](server/client_stream.go)

---

### 4. Bidirectional Streaming

```
Client                                  Server
  |-- HelloRequest (name[0]) ---------->|
  |<-- HelloResponse (name[0]) ---------|
  |-- HelloRequest (name[1]) ---------->|
  |<-- HelloResponse (name[1]) ---------|
  |-- CloseSend() --------------------->|
  |<-- EOF -----------------------------|
```

- Client launches a goroutine to receive responses concurrently.
- The send loop iterates names, sends a `HelloRequest`, then sleeps 2 seconds.
- Server echoes a `HelloResponse` for each received request immediately.
- After all names are sent, the client calls `stream.CloseSend()` and waits on a channel (`waitc`) for the receive goroutine to drain and signal `io.EOF`.

**Files:** [`client/bi_stream.go`](client/bi_stream.go), [`server/bi_stream.go`](server/bi_stream.go)

---

## Prerequisites

- Go 1.21 or later — [go.dev/dl](https://go.dev/dl/)
- `protoc` — Protocol Buffers compiler — [github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Both plugins must be on your `PATH`.

---

## Code Generation

Run this from the project root whenever `greet.proto` changes:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/greet.proto
```

This overwrites `proto/greet.pb.go` and `proto/greet_grpc.pb.go`. The generated files are committed to the repository so you only need to re-run this if you modify the `.proto` file.

---

## Running

**1. Install dependencies**

```bash
go mod tidy
```

**2. Start the server**

```bash
go run ./server
```

The server listens on `localhost:8080`.

**3. Run the client** (in a separate terminal)

```bash
go run ./client
```

By default, `client/main.go` calls `callSayHelloBiDirectionalStreaming`. To exercise the other patterns, comment out that line and uncomment the corresponding call:

```go
// callSayHello(client)
// callSayHelloServerStream(client, names)
// callSayHelloClientStream(client, names)
callSayHelloBiDirectionalStreaming(client, names)
```

---

## Dependencies

| Module                                      | Version  | Role                           |
| ------------------------------------------- | -------- | ------------------------------ |
| `google.golang.org/grpc`                    | v1.84.0  | gRPC runtime, HTTP/2 transport |
| `google.golang.org/protobuf`                | v1.36.12 | Protobuf serialization         |
| `golang.org/x/net`                          | v0.57.0  | Indirect (HTTP/2 internals)    |
| `golang.org/x/sys`                          | v0.47.0  | Indirect (OS syscalls)         |
| `golang.org/x/text`                         | v0.40.0  | Indirect (text encoding)       |
| `google.golang.org/genproto/googleapis/rpc` | latest   | Indirect (status codes)        |

---

Built by [Jeet Das](https://github.com/JeetDas5)
