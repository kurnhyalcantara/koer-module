# koer-module

A shared Go module providing reusable infrastructure packages for Koer platform backend microservices. This is a **library** — not a runnable service — consumed as a Go module dependency.

## Packages

| Package | Description |
|---------|-------------|
| `pkg/config` | `.env`-based configuration loading, centralized config structs, and `.env` template generation |
| `pkg/connection` | Clients: MySQL, Redis, Kafka (producer & consumer), MinIO, Firebase, REST |
| `pkg/server` | Server setups: HTTP (Fiber), gRPC, and combined HTTP+gRPC |
| `pkg/jwt` | JWT token creation and validation |
| `pkg/logger` | Structured logging via zerolog |
| `pkg/tracing` | OpenTelemetry/APM distributed tracing |
| `pkg/apiresponse` | Standardized API response models and helpers |
| `pkg/utils` | Utility helpers (slices, random generation) |

## Requirements

- Go 1.25+

## Installation

```bash
go get github.com/kurnhyalcantara/koer-module
```

## Configuration

All config structs are centralized in `pkg/config`. Services import and compose only what they need.

### Available Config Structs

| Struct | Environment Variables |
|--------|----------------------|
| `MySQLConfig` | `MYSQL_DSN`, `MYSQL_MAX_OPEN_CONNS`, `MYSQL_MAX_IDLE_CONNS`, `MYSQL_CONN_MAX_LIFETIME` |
| `RedisConfig` | `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB` |
| `KafkaProducerConfig` | `KAFKA_PRODUCER_BROKERS`, `KAFKA_PRODUCER_TOPIC` |
| `KafkaConsumerConfig` | `KAFKA_CONSUMER_BROKERS`, `KAFKA_CONSUMER_TOPIC`, `KAFKA_CONSUMER_GROUP_ID` |
| `FirebaseConfig` | `FIREBASE_CREDENTIALS_FILE`, `FIREBASE_PROJECT_ID` |
| `MinIOConfig` | `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY_ID`, `MINIO_SECRET_ACCESS_KEY`, `MINIO_USE_SSL`, `MINIO_REGION` |
| `RESTClientConfig` | `REST_BASE_URL`, `REST_TIMEOUT` |
| `JWTConfig` | `JWT_SECRET_KEY`, `JWT_ACCESS_EXPIRY`, `JWT_REFRESH_EXPIRY`, `JWT_ISSUER` |
| `LoggerConfig` | `LOGGER_LEVEL`, `LOGGER_PRETTY` |
| `TracingConfig` | `TRACING_SERVICE_NAME`, `TRACING_ENDPOINT`, `TRACING_ENABLED` |
| `HTTPConfig` | `HTTP_PORT`, `HTTP_READ_TIMEOUT`, `HTTP_WRITE_TIMEOUT` |
| `GRPCConfig` | `GRPC_PORT` |

Use `AppConfig` to load all of them at once, or compose a custom struct from the ones you need.

### Loading Config

```go
import "github.com/kurnhyalcantara/koer-module/pkg/config"

// Load from .env file in current directory
var cfg config.AppConfig
if err := config.Load("", &cfg); err != nil {
    log.Fatal(err)
}

// Or load a specific file
if err := config.Load("/etc/myservice/.env", &cfg); err != nil {
    log.Fatal(err)
}
```

**Compose only what your service needs:**

```go
type MyServiceConfig struct {
    DB     config.MySQLConfig
    Cache  config.RedisConfig
    JWT    config.JWTConfig
    Logger config.LoggerConfig
}

var cfg MyServiceConfig
config.Load("", &cfg)
```

### Generating a .env Template

`GenerateConfig` writes a `.env` template pre-filled with default values. Pass your own config structs to include service-specific variables alongside the mandatory ones.

```go
type ServiceConfig struct {
    AppName    string `env:"APP_NAME"  envDefault:"my-service"`
    FeatureXOn bool   `env:"FEATURE_X" envDefault:"false"`
}

// generates .env with all library sections + your ServiceConfig section
if err := config.GenerateConfig(".env", ServiceConfig{}); err != nil {
    log.Fatal(err)
}
```

The section name in the generated file is derived from your struct type name.

## Usage Examples

### HTTP + gRPC Server

```go
var cfg config.AppConfig
config.Load("", &cfg)

httpServer := server.NewHTTPServer(cfg.HTTP)
grpcServer := server.NewGRPCServer(cfg.GRPC)

// or combined
combined := server.NewCombinedServer(config.CombinedConfig{
    HTTP: cfg.HTTP,
    GRPC: cfg.GRPC,
})
combined.Start()
```

### Database & Cache

```go
db, err := connection.NewMySQLClient(cfg.MySQL)
redisClient := connection.NewRedisClient(cfg.Redis)
```

### JWT

```go
jwtManager := jwt.NewManager(cfg.JWT)

accessToken, err := jwtManager.GenerateAccessToken(userID, role)
claims, err := jwtManager.ValidateToken(tokenStr)
```

### Logger

```go
log := logger.New(cfg.Logger)
log.Info().Str("key", "value").Msg("service started")
```

### Tracing

```go
provider, err := tracing.NewProvider(ctx, cfg.Tracing)
defer provider.Shutdown(ctx)
```

## Development

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Tidy dependencies
go mod tidy
```

## License

MIT
