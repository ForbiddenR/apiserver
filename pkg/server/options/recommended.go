package options

import "github.com/ForbiddenR/apiserver/pkg/server"

// RecommendedOptions contains the recommended options for running an API server.
// If you add something to this list, it should be in a logical grouping.
// Each of them can be nil to leave the feature unconfigured on ApplyTo.
type RecommendedOptions struct {
	Serving *ServingOptions
	CoreAPI *CoreAPIOptions
}

func NewRecommendedOptions() *RecommendedOptions {

	return &RecommendedOptions{
		CoreAPI: NewCoreAPIOptions(),
		Serving: NewServingOptions(),
	}
}

func (o *RecommendedOptions) ApplyTo(config *server.RecommendedConfig) error {
	if err := o.CoreAPI.ApplyTo(config); err != nil {
		return err
	}
	if err := o.Serving.ApplyTo(&config.Config.Serving); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	errors := []error{}
	errors = append(errors, o.CoreAPI.Validate()...)

	return errors
}