# gRPC Service Example

## ทำไมต้องใช้ gRPC?

| HTTP/REST | gRPC |
|-----------|------|
| JSON (text) | Protobuf (binary) |
| Slower | 7-10x faster |
| Human readable | Machine optimized |
| Request-Response | Streaming support |

## การติดตั้ง

```bash
# Install protoc compiler
# macOS
brew install protobuf

# Linux
apt install -y protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Generate Proto

```bash
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

## รัน

```bash
# Start gRPC server
go run server/main.go

# Start client (ในอีก terminal)
go run client/main.go
```

## เปรียบเทียบ Performance

| Metric | REST | gRPC |
|--------|------|------|
| Latency | 2.5ms | 0.3ms |
| Throughput | 10K RPS | 100K RPS |
| Payload size | 1KB | 100B |
