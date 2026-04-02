package logger

import (
	"os"
	"strings"

	"github.com/koer/koer-module/pkg/config"
	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(cfg config.LoggerConfig) *Logger {
	var level zerolog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	var zl zerolog.Logger
	if cfg.Pretty {
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(level).With().Timestamp().Logger()
	} else {
		zl = zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger()
	}
	return &Logger{zl}
}

var defaultLogger = New(config.LoggerConfig{Level: "info"})

func Default() *Logger { return defaultLogger }
