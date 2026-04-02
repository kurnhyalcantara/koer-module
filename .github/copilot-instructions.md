# Copilot Instructions — koer-module

## Overview

`koer-module` is a shared Go module providing reusable infrastructure packages for Koer platform backend microservices. It is a **library** — not a runnable service — consumed as a Go module dependency by other microservices.

## Package Architecture

```
pkg/
  connection/   – Clients: MySQL, Redis, Kafka (producer & consumer), MinIO, Firebase, REST
  server/       – Server setups: HTTP (Fiber), gRPC, and combined HTTP+gRPC
  jwt/          – JWT token creation and validation
  config/       – Configuration loading
  logger/       – Structured logging via zerolog
  tracing/      – OpenTelemetry/APM distributed tracing
  apiresponse/  – Standardized API response models, helpers, and protobuf definitions
  utils/        – Utility helpers (slices, random generation)
```

## Build & Test Commands

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run tests for a single package
go test ./pkg/jwt/...

# Run a single test by name
go test ./pkg/jwt/... -run TestTokenValidation

# Run tests with race detector
go test -race ./...

# Lint (assumes golangci-lint is installed)
golangci-lint run ./...

# Tidy dependencies
go mod tidy
```

## Key Conventions

### This is a library — no `main` package
All packages live under `pkg/`. No `cmd/` or `main.go`. Consumers import specific sub-packages.

### Connection clients
Each client in `pkg/connection` should expose a constructor (e.g., `NewMySQLClient`, `NewRedisClient`) that accepts a config struct or options. Clients must be safe for concurrent use.

### Server setup
`pkg/server` wires together Fiber (HTTP) and gRPC listeners. Combined HTTP+gRPC setups use a single port via protocol detection or separate ports — keep both modes supported.

### Configuration
`pkg/config` handles loading; config structs for each package should be defined close to the package that uses them (not centralized), and composed by the consuming service.

### Logging
Use `pkg/logger` (zerolog) throughout. Do not use `fmt.Println` or `log.Printf` in library code — return errors or use the structured logger.

### API responses
`pkg/apiresponse` defines the canonical response envelope used across all microservices. Protobuf definitions live alongside the Go helpers here. Any change to the response schema affects all consumers — treat it carefully.

### Error handling
Return errors explicitly; avoid panics in library code. Wrap errors with context using `fmt.Errorf("...: %w", err)`.

### Tracing
`pkg/tracing` wraps OpenTelemetry setup. Span creation should follow OTel conventions; use `trace.SpanFromContext(ctx)` to propagate spans rather than creating new root spans mid-call.

### Module path
When adding new packages, follow the existing module path (check `go.mod`). Keep packages focused — one responsibility per `pkg/` subdirectory.
