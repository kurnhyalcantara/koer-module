package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// GenerateConfig writes a .env template file derived from the provided
// config structs. It inspects each struct's fields: nested struct fields
// (without an `env` tag) are treated as sections, and fields with an
// `env` tag are written as environment variable declarations with their
// default values. If filename is empty, ".env" is used.
//
// Consumer services define their own config struct composing only the
// infra configs they need:
//
//	type Config struct {
//	    GRPC   config.GRPCConfig
//	    MySQL  config.MySQLConfig
//	    Logger config.LoggerConfig
//	    AppName string `env:"APP_NAME" envDefault:"my-service"`
//	}
//	config.GenerateConfig(".env", Config{})
func GenerateConfig(filename string, configs ...any) error {
	if filename == "" {
		filename = ".env"
	}

	var sb strings.Builder
	sb.WriteString("# Auto-generated .env configuration\n")
	sb.WriteString("# Edit values before running your service.\n\n")

	for _, cfg := range configs {
		t := reflect.TypeOf(cfg)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		writeConfigStruct(&sb, t)
	}

	return os.WriteFile(filename, []byte(sb.String()), 0o644)
}

// writeConfigStruct iterates over a struct's fields. Nested struct fields
// without an `env` tag are emitted as named sections; fields with an
// `env` tag are written directly as env variable declarations.
func writeConfigStruct(sb *strings.Builder, t reflect.Type) {
	for i := range t.NumField() {
		field := t.Field(i)

		// Nested struct without its own env tag → section.
		if field.Type.Kind() == reflect.Struct && field.Tag.Get("env") == "" {
			sb.WriteString(fmt.Sprintf("# --- %s ---\n", field.Name))
			writeEnvFields(sb, field.Type)
			sb.WriteString("\n")
			continue
		}

		varName := field.Tag.Get("env")
		if varName == "" {
			continue
		}

		defaultVal := field.Tag.Get("envDefault")
		sb.WriteString(fmt.Sprintf("%s=%s\n", varName, defaultVal))
	}
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
