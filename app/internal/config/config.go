package config

import (
	pkgconfig "github.com/alfattd/category-service/internal/pkg/config"
)

type Config struct {
	pkgconfig.Base

	RabbitMQUrl string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
}

// internal/config/config.go
func Load() *Config {
	return &Config{
		Base: pkgconfig.LoadBase(), // load generic sekaligus

		RabbitMQUrl: pkgconfig.Env("RABBITMQ_URL", ""),
		DBHost:      pkgconfig.Env("DB_HOST", ""),
		DBPort:      pkgconfig.Env("DB_PORT", "5432"),
		DBName:      pkgconfig.Env("DB_NAME", ""),
		DBUser:      pkgconfig.Env("DB_USER", ""),
		DBPassword:  pkgconfig.Env("DB_PASSWORD", ""),
		DBSSLMode:   pkgconfig.Env("DB_SSLMODE", "disable"),
	}
}

func (c *Config) Validate() error {
	if err := c.ValidateBase(); err != nil { // validate generic dulu
		return err
	}

	fields := []struct{ value, name string }{
		{c.DBHost, "DB_HOST"},
		{c.DBPort, "DB_PORT"},
		{c.DBUser, "DB_USER"},
		{c.DBPassword, "DB_PASSWORD"},
		{c.DBName, "DB_NAME"},
		{c.RabbitMQUrl, "RABBITMQ_URL"},
	}
	for _, f := range fields {
		if err := pkgconfig.Required(f.value, f.name); err != nil {
			return err
		}
	}

	return nil
}
