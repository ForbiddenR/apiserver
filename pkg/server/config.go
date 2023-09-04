package server

import "github.com/ForbiddenR/apiserver/pkg/server/healthz"

type Config struct {
	// The default set of healthz checks. There might be more added via AddHealthChecks dynamically.
	HealthzChecks []healthz.HealthzChecker
	// The default set of livez checks. There might be more added via AddHealthChecks dynamically.
	LivezChecks []healthz.HealthzChecker
	// The default set of readyz-only checks. There might be more added via AddReadyzChecks dynamically.
	ReadyzChecks []healthz.HealthzChecker
	// If specified, long running requests such as watch will be allocated a random timeout between this value, and
	// twice this value. Note that it is up to the request handlers to ignore or honor this timeout. In seconds.
	MinRequestTimeout int
}

type RecommendedConfig struct {
	Config
}

// NewConfig returns a Config struct with default values.
func NewConfig() *Config {
	defaultHeathChecks := []healthz.HealthzChecker{healthz.LivenessHealthz, healthz.ReadinessHealthz}

	return &Config{
		HealthzChecks: append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		LivezChecks:   append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		ReadyzChecks:  append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		MinRequestTimeout: 1800,
	}
}

func NewRecommendedConfig() *RecommendedConfig {
	return &RecommendedConfig{
		Config: *NewConfig(),
	}
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
