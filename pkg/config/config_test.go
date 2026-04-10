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

func TestGenerateConfig_ComposedStruct(t *testing.T) {
	type ServiceConfig struct {
		GRPC          GRPCConfig
		MySQL         MySQLConfig
		Logger        LoggerConfig
		PublicMethods []string `env:"PUBLIC_GRPC_METHODS" envSeparator:"," envDefault:""`
	}

	dir := t.TempDir()
	filename := dir + "/test-composed.env"

	if err := GenerateConfig(filename, ServiceConfig{}); err != nil {
		t.Fatalf("GenerateConfig: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}
	content := string(data)

	// section headers derived from field names
	for _, section := range []string{"# --- GRPC ---", "# --- MySQL ---", "# --- Logger ---"} {
		if !strings.Contains(content, section) {
			t.Errorf("missing section header %q", section)
		}
	}
	// entries from composed structs
	if !strings.Contains(content, "GRPC_PORT=9090") {
		t.Error("missing GRPC_PORT entry")
	}
	if !strings.Contains(content, "MYSQL_DSN=") {
		t.Error("missing MYSQL_DSN entry")
	}
	if !strings.Contains(content, "LOGGER_LEVEL=info") {
		t.Error("missing LOGGER_LEVEL entry")
	}
	// top-level env-tagged field
	if !strings.Contains(content, "PUBLIC_GRPC_METHODS=") {
		t.Error("missing PUBLIC_GRPC_METHODS entry")
	}
	// should NOT contain configs not in the composed struct
	if strings.Contains(content, "HTTP_PORT") {
		t.Error("unexpected HTTP_PORT entry — should not be generated")
	}
}

func TestGenerateConfig_AppConfig(t *testing.T) {
	dir := t.TempDir()
	filename := dir + "/test.env"

	if err := GenerateConfig(filename, AppConfig{}); err != nil {
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
