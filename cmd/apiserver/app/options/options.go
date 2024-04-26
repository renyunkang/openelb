package options

import (
	server "github.com/openelb/openelb/pkg/server/options"
	cliflag "k8s.io/component-base/cli/flag"
)

type OpenELBApiServerOptions struct {
	HTTPOptions *server.Options
}

func NewOpenELBApiServerOptions() *OpenELBApiServerOptions {
	return &OpenELBApiServerOptions{
		HTTPOptions: server.NewOptions(),
	}
}

func (s *OpenELBApiServerOptions) Validate() []error {
	var errs []error
	return errs
}

func (s *OpenELBApiServerOptions) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}

	s.HTTPOptions.AddFlags(fss.FlagSet("http"))

	return fss
}
