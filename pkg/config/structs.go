package config

import "time"

// MySQLConfig holds MySQL connection configuration.
type MySQLConfig struct {
	DSN             string        `env:"MYSQL_DSN"`
	MaxOpenConns    int           `env:"MYSQL_MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns    int           `env:"MYSQL_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"MYSQL_CONN_MAX_LIFETIME" envDefault:"1h"`
}

// PostgreSQLConfig holds PostgreSQL connection configuration.
type PostgreSQLConfig struct {
	DSN             string        `env:"POSTGRES_DSN"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"1h"`
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

// KafkaProducerConfig holds Kafka producer configuration.
type KafkaProducerConfig struct {
	Brokers []string `env:"KAFKA_PRODUCER_BROKERS" envDefault:"localhost:9092" envSeparator:","`
	Topic   string   `env:"KAFKA_PRODUCER_TOPIC"`
}

// KafkaConsumerConfig holds Kafka consumer configuration.
type KafkaConsumerConfig struct {
	Brokers []string `env:"KAFKA_CONSUMER_BROKERS" envDefault:"localhost:9092" envSeparator:","`
	Topic   string   `env:"KAFKA_CONSUMER_TOPIC"`
	GroupID string   `env:"KAFKA_CONSUMER_GROUP_ID"`
}

// FirebaseConfig holds Firebase connection configuration.
type FirebaseConfig struct {
	CredentialsFile string `env:"FIREBASE_CREDENTIALS_FILE"`
	ProjectID       string `env:"FIREBASE_PROJECT_ID"`
}

// MinIOConfig holds MinIO connection configuration.
type MinIOConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	Region          string `env:"MINIO_REGION" envDefault:"us-east-1"`
}

// RESTClientConfig holds REST client configuration.
// Headers cannot be set via environment variables; assign programmatically.
type RESTClientConfig struct {
	BaseURL string            `env:"REST_BASE_URL"`
	Timeout time.Duration     `env:"REST_TIMEOUT" envDefault:"30s"`
	Headers map[string]string // set programmatically, not via env
}

// JWTConfig holds JWT token configuration.
type JWTConfig struct {
	SecretKey     string        `env:"JWT_SECRET_KEY"`
	AccessExpiry  time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"168h"`
	Issuer        string        `env:"JWT_ISSUER"`
}

// LoggerConfig holds structured logger configuration.
type LoggerConfig struct {
	Level  string `env:"LOGGER_LEVEL" envDefault:"info"`
	Pretty bool   `env:"LOGGER_PRETTY" envDefault:"false"`
}

// TracingConfig holds OpenTelemetry tracing configuration.
type TracingConfig struct {
	ServiceName string `env:"TRACING_SERVICE_NAME"`
	Endpoint    string `env:"TRACING_ENDPOINT"`
	Enabled     bool   `env:"TRACING_ENABLED" envDefault:"false"`
}

// HTTPConfig holds HTTP server configuration.
type HTTPConfig struct {
	Port         int           `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
}

// GRPCConfig holds gRPC server configuration.
type GRPCConfig struct {
	Port int `env:"GRPC_PORT" envDefault:"9090"`
}

// CombinedConfig holds combined HTTP and gRPC server configuration.
type CombinedConfig struct {
	HTTP HTTPConfig
	GRPC GRPCConfig
}

// AppConfig is the master configuration struct containing all service configs.
// Pass a pointer to this (or any subset) to config.Load to populate from env.
type AppConfig struct {
	MySQL         MySQLConfig
	Redis         RedisConfig
	KafkaProducer KafkaProducerConfig
	KafkaConsumer KafkaConsumerConfig
	Firebase      FirebaseConfig
	MinIO         MinIOConfig
	RESTClient    RESTClientConfig
	JWT           JWTConfig
	Logger        LoggerConfig
	Tracing       TracingConfig
	HTTP          HTTPConfig
	GRPC          GRPCConfig
}
