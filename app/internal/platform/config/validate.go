package config

func (c *Config) Validate() error {

	if err := required(c.AppPort, "APP_PORT"); err != nil {
		return err
	}

	if err := required(c.ServiceName, "SERVICE_NAME"); err != nil {
		return err
	}

	if err := required(c.ServiceVersion, "SERVICE_VERSION"); err != nil {
		return err
	}

	if err := required(c.RabbitMQUrl, "RABBITMQ_URL"); err != nil {
		return err
	}

	if c.ServiceVersion != "dev" {

		if err := required(c.DBHost, "DB_HOST"); err != nil {
			return err
		}

		if err := required(c.DBPort, "DB_PORT"); err != nil {
			return err
		}

		if err := required(c.DBUser, "DB_USER"); err != nil {
			return err
		}

		if err := required(c.DBPassword, "DB_PASSWORD"); err != nil {
			return err
		}

		if err := required(c.DBName, "DB_NAME"); err != nil {
			return err
		}

		if err := required(c.DBSSLMode, "DB_SSLMODE"); err != nil {
			return err
		}
	}

	return nil
}
