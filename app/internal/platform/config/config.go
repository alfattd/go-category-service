package config

import "os"

type Config struct {
	AppPort        string
	ServiceName    string
	ServiceVersion string

	RabbitMQUrl string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {
	return &Config{
		AppPort:        os.Getenv("APP_PORT"),
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceVersion: os.Getenv("SERVICE_VERSION"),

		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

func (c *Config) DBUrl() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" dbname=" + c.DBName +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" sslmode=" + c.DBSSLMode
}
