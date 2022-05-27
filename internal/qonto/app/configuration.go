package app

// Configuration holds application configuration
type Configuration struct {
	ListenAddress string
	DB            struct {
		Address  string
		User     string
		Password string
		Name     string
	}
}

// ConfigurationFromEnv builds configuration of the app based on env
func ConfigurationFromEnv(envGetter EnvGetterFunc) (*Configuration, error) {
	config := Configuration{}
	listenAddr := envGetter("QONTO_APP_LISTEN_ADDRESS")
	if listenAddr == "" {
		listenAddr = "127.0.0.1:8080"
	}

	config.ListenAddress = listenAddr
	config.DB.Address = envGetter("QONTO_DB_ADDRESS")
	config.DB.Name = envGetter("QONTO_DB_NAME")
	config.DB.Password = envGetter("QONTO_DB_PASSWORD")
	config.DB.User = envGetter("QONTO_DB_USER")

	return &config, nil
}
