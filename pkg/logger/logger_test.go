package logger

import (
	"testing"

	"github.com/koer/koer-module/pkg/config"
)

func TestNew(t *testing.T) {
	l := New(config.LoggerConfig{Level: "debug", Pretty: false})
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestDefault(t *testing.T) {
	l := Default()
	if l == nil {
		t.Fatal("expected non-nil default logger")
	}
}
