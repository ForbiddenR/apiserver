package server

type Config struct {
}

type RecommendedConfig struct {
	Config
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

func (c completedConfig) New() (*GenericAPIServer, error) {
	apiServerHandler := NewAPIServerHandler()

	s := &GenericAPIServer{
		Handler: apiServerHandler,
	}

	installAPI(s, c.Config)

	return s, nil
}

func installAPI(s *GenericAPIServer, c *Config) {
}
