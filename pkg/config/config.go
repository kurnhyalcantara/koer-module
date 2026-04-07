package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Load reads environment variables from an optional .env file and parses
// them into v using struct field `env` tags. If envFile is empty, it
// attempts to load ".env" from the current directory. When the .env file
// does not exist, it falls back to reading directly from os.Getenv using
// the struct's `env` tags.
func Load(envFile string, v any) error {
	var fileNotFound bool
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fileNotFound = true
			} else {
				return fmt.Errorf("loading env file %q: %w", envFile, err)
			}
		}
	} else {
		if err := godotenv.Load(); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fileNotFound = true
			}
		}
	}

	if fileNotFound {
		return loadFromOsEnv(v)
	}

	if err := env.Parse(v); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	return nil
}

// loadFromOsEnv populates struct fields in v by reading environment
// variables directly via os.Getenv, using the `env` struct tag as key.
func loadFromOsEnv(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("expected non-nil pointer to struct, got %T", v)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got pointer to %s", rv.Kind())
	}
	return parseStructFields(rv)
}

func parseStructFields(rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// Recurse into embedded or nested structs.
		if fieldVal.Kind() == reflect.Struct && field.Tag.Get("env") == "" {
			if err := parseStructFields(fieldVal); err != nil {
				return err
			}
			continue
		}

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		val := os.Getenv(envKey)
		if val == "" {
			val = field.Tag.Get("envDefault")
		}
		if val == "" {
			continue
		}

		if err := setField(fieldVal, val, field.Tag.Get("envSeparator")); err != nil {
			return fmt.Errorf("setting field %s: %w", field.Name, err)
		}
	}
	return nil
}

func setField(fieldVal reflect.Value, val, separator string) error {
	// Handle time.Duration before the kind switch (it's int64 underneath).
	if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(d))
		return nil
	}

	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		fieldVal.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		fieldVal.SetUint(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fieldVal.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		fieldVal.SetFloat(f)
	case reflect.Slice:
		if fieldVal.Type().Elem().Kind() == reflect.String {
			sep := ","
			if separator != "" {
				sep = separator
			}
			parts := strings.Split(val, sep)
			fieldVal.Set(reflect.ValueOf(parts))
		}
	}
	return nil
}
