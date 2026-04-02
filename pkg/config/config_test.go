package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_NoFile(t *testing.T) {
	var cfg AppConfig
	if err := Load("", &cfg); err != nil {
		t.Fatalf("Load with no file should not error: %v", err)
	}
}

func TestLoad_WithEnvVars(t *testing.T) {
	t.Setenv("HTTP_PORT", "9999")
	t.Setenv("LOGGER_LEVEL", "debug")

	var cfg AppConfig
	if err := Load("", &cfg); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTP.Port != 9999 {
		t.Errorf("expected HTTP.Port=9999, got %d", cfg.HTTP.Port)
	}
	if cfg.Logger.Level != "debug" {
		t.Errorf("expected Logger.Level=debug, got %s", cfg.Logger.Level)
	}
}

func TestLoad_WithEnvFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("GRPC_PORT=7777\n")
	f.Close()

	var cfg AppConfig
	if err := Load(f.Name(), &cfg); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.GRPC.Port != 7777 {
		t.Errorf("expected GRPC.Port=7777, got %d", cfg.GRPC.Port)
	}
}

func TestGenerateConfig_WithExtras(t *testing.T) {
	type ServiceConfig struct {
		AppName string `env:"APP_NAME" envDefault:"my-service"`
		Debug   bool   `env:"APP_DEBUG" envDefault:"false"`
	}

	dir := t.TempDir()
	filename := dir + "/test-extras.env"

	if err := GenerateConfig(filename, ServiceConfig{}); err != nil {
		t.Fatalf("GenerateConfig: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}
	content := string(data)

	// mandatory entry still present
	if !strings.Contains(content, "HTTP_PORT=8080") {
		t.Error("missing mandatory HTTP_PORT entry")
	}
	// consumer-defined entries present
	if !strings.Contains(content, "APP_NAME=my-service") {
		t.Error("missing APP_NAME from extra config")
	}
	if !strings.Contains(content, "APP_DEBUG=false") {
		t.Error("missing APP_DEBUG from extra config")
	}
	// section header derived from type name
	if !strings.Contains(content, "# --- ServiceConfig ---") {
		t.Error("missing ServiceConfig section header")
	}
}

func TestGenerateConfig(t *testing.T) {
	dir := t.TempDir()
	filename := dir + "/test.env"

	if err := GenerateConfig(filename); err != nil {
		t.Fatalf("GenerateConfig: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}
	content := string(data)

	expected := []string{
		"MYSQL_DSN=",
		"REDIS_ADDR=localhost:6379",
		"HTTP_PORT=8080",
		"GRPC_PORT=9090",
		"JWT_SECRET_KEY=",
		"LOGGER_LEVEL=info",
		"TRACING_ENABLED=false",
	}
	for _, want := range expected {
		if !strings.Contains(content, want) {
			t.Errorf("generated .env missing %q", want)
		}
	}
}
