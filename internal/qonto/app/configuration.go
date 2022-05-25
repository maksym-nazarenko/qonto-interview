package app

// Configuration holds application configuration
type Configuration struct {
	DB struct {
		Address  string
		User     string
		Password string
		Name     string
	}
}

// ConfigurationFromEnv builds configuration of the app based on env
func ConfigurationFromEnv(envGetter EnvGetterFunc) (*Configuration, error) {
	config := Configuration{}

	config.DB.Address = envGetter("QONTO_DB_ADDRESS")
	config.DB.Name = envGetter("QONTO_DB_NAME")
	config.DB.Password = envGetter("QONTO_DB_PASSWORD")
	config.DB.User = envGetter("QONTO_DB_USER")

	return &config, nil
}
