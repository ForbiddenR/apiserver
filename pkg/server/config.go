package server

import (
	"net"
	"time"

	"github.com/ForbiddenR/apiserver/pkg/server/healthz"
)

type Config struct {
	// Serving is required to serve http
	Serving *ServingInfo
	// The default set of healthz checks. There might be more added via AddHealthChecks dynamically.
	HealthzChecks []healthz.HealthzChecker
	// The default set of livez checks. There might be more added via AddHealthChecks dynamically.
	LivezChecks []healthz.HealthzChecker
	// The default set of readyz-only checks. There might be more added via AddReadyzChecks dynamically.
	ReadyzChecks []healthz.HealthzChecker
	// If specified, all requests except those which match the LongRunningFunc predicate will timeout
	// after this duration.
	RequestTimeout time.Duration
	// If specified, long running requests such as watch will be allocated a random timeout between this value, and
	// twice this value. Note that it is up to the request handlers to ignore or honor this timeout. In seconds.
	MinRequestTimeout int
}

// Complete fills in any fields not set that are required to have valid data and can be drived
// from othe fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{c}}
}

type RecommendedConfig struct {
	Config
}

type ServingInfo struct {
	// Listener is the secure server network listener.
	Listener net.Listener
}

// NewConfig returns a Config struct with default values.
func NewConfig() *Config {
	defaultHeathChecks := []healthz.HealthzChecker{healthz.LivenessHealthz, healthz.ReadinessHealthz}

	return &Config{
		HealthzChecks:     append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		LivezChecks:       append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		ReadyzChecks:      append([]healthz.HealthzChecker{}, defaultHeathChecks...),
		RequestTimeout:    time.Duration(60) * time.Second,
		MinRequestTimeout: 1800,
	}
}

func NewRecommendedConfig() *RecommendedConfig {
	return &RecommendedConfig{
		Config: *NewConfig(),
	}
}

// Complete fills in any fields not set that are required to have valid data and can be drived
// from othe fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *RecommendedConfig) Complete() CompletedConfig {
	return c.Config.Complete()
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging.
func (c completedConfig) New(name string) (*GenericAPIServer, error) {
	apiServerHandler := NewAPIServerHandler()

	s := &GenericAPIServer{
		Handler: apiServerHandler,

		minRequestTimeout: time.Duration(c.MinRequestTimeout) * time.Second,
		ShutdownTimeout:   c.RequestTimeout,
		ServingInfo:       c.Serving,

		healthzChecks: c.HealthzChecks,
		livezChecks:   c.LivezChecks,
		readyzChecks:  c.ReadyzChecks,
	}

	installAPI(s, c.Config)

	return s, nil
}

func installAPI(s *GenericAPIServer, c *Config) {
}
