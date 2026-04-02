package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// GenerateConfig writes a .env template file derived from all mandatory
// config structs defined in this package, followed by any extra config
// structs provided by the consumer service. Each variable is written with
// its default value, or left empty if no default is defined.
// If filename is empty, ".env" is used.
//
// Example:
//
//	type ServiceConfig struct {
//	    AppName string `env:"APP_NAME" envDefault:"my-service"`
//	    Debug   bool   `env:"APP_DEBUG" envDefault:"false"`
//	}
//	config.GenerateConfig(".env", ServiceConfig{})
func GenerateConfig(filename string, extras ...any) error {
	if filename == "" {
		filename = ".env"
	}

	var sb strings.Builder
	sb.WriteString("# Auto-generated .env configuration\n")
	sb.WriteString("# Edit values before running your service.\n\n")

	sections := []struct {
		name string
		cfg  any
	}{
		{"MySQL", MySQLConfig{}},
		{"Redis", RedisConfig{}},
		{"Kafka Producer", KafkaProducerConfig{}},
		{"Kafka Consumer", KafkaConsumerConfig{}},
		{"Firebase", FirebaseConfig{}},
		{"MinIO", MinIOConfig{}},
		{"REST Client", RESTClientConfig{}},
		{"JWT", JWTConfig{}},
		{"Logger", LoggerConfig{}},
		{"Tracing", TracingConfig{}},
		{"HTTP Server", HTTPConfig{}},
		{"gRPC Server", GRPCConfig{}},
	}

	for _, s := range sections {
		sb.WriteString(fmt.Sprintf("# --- %s ---\n", s.name))
		writeEnvFields(&sb, reflect.TypeOf(s.cfg))
		sb.WriteString("\n")
	}

	for _, extra := range extras {
		t := reflect.TypeOf(extra)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		sb.WriteString(fmt.Sprintf("# --- %s ---\n", t.Name()))
		writeEnvFields(&sb, t)
		sb.WriteString("\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0o644)
}

// writeEnvFields iterates over struct fields and writes each env variable
// declaration. Fields without an `env` tag (e.g. map fields set
// programmatically) are silently skipped.
func writeEnvFields(sb *strings.Builder, t reflect.Type) {
	for i := range t.NumField() {
		field := t.Field(i)

		varName := field.Tag.Get("env")
		if varName == "" {
			continue
		}

		defaultVal := field.Tag.Get("envDefault")
		sb.WriteString(fmt.Sprintf("%s=%s\n", varName, defaultVal))
	}
}
