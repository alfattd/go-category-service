package config

type Base struct {
	AppPort        string
	ServiceName    string
	ServiceVersion string
}

func LoadBase() Base {
	return Base{
		AppPort:        Env("APP_PORT", "80"),
		ServiceName:    Env("SERVICE_NAME", ""),
		ServiceVersion: Env("SERVICE_VERSION", ""),
	}
}

func (b Base) ValidateBase() error {
	fields := []struct{ value, name string }{
		{b.AppPort, "APP_PORT"},
		{b.ServiceName, "SERVICE_NAME"},
		{b.ServiceVersion, "SERVICE_VERSION"},
	}
	for _, f := range fields {
		if err := Required(f.value, f.name); err != nil {
			return err
		}
	}
	return nil
}
