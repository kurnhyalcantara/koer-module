package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Load reads environment variables from an optional .env file and parses
// them into v using struct field `env` tags. If envFile is empty, it
// attempts to load ".env" from the current directory. A missing .env file
// is silently ignored — variables already present in the environment are
// always used regardless.
func Load(envFile string, v any) error {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("loading env file %q: %w", envFile, err)
		}
	} else {
		_ = godotenv.Load() // silently ignore missing .env
	}
	if err := env.Parse(v); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	return nil
}
