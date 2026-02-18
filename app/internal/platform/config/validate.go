package config

import "fmt"

func (c *Config) Validate() error {
	fields := []struct{ value, name string }{
		{c.AppPort, "APP_PORT"},
		{c.ServiceName, "SERVICE_NAME"},
		{c.ServiceVersion, "SERVICE_VERSION"},
		{c.DBHost, "DB_HOST"},
		{c.DBPort, "DB_PORT"},
		{c.DBUser, "DB_USER"},
		{c.DBPassword, "DB_PASSWORD"},
		{c.DBName, "DB_NAME"},
		{c.RabbitMQUrl, "RABBITMQ_URL"},
	}

	for _, f := range fields {
		if err := required(f.value, f.name); err != nil {
			return err
		}
	}

	return nil
}

func required(value, name string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}
