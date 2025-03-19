package options

import "github.com/ForbiddenR/apiserver/pkg/server"

type CoreAPIOptions struct {
	CoreAPIPath string
}

func NewCoreAPIOptions() *CoreAPIOptions {
	return &CoreAPIOptions{}
}

func (o *CoreAPIOptions) ApplyTo(config *server.RecommendedConfig) error {
	return nil
}

func (o *CoreAPIOptions) Validate() []error {
	return nil
}
