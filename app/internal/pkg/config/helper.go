package config

import (
	"fmt"
	"os"
)

func Env(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Required(value, name string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}
