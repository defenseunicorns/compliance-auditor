package api

import (
	"github.com/defenseunicorns/lula/src/types"
)

type ApiDomain struct {
	// Spec is the specification of the API requests
	Spec *ApiSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (a ApiDomain) GetResources() (types.DomainResources, error) {
	return MakeRequests(a.Spec.Requests)
}

func (a ApiDomain) IsExecutable() bool {
	// Domain is not currently executable
	return false
}

type ApiSpec struct {
	Requests []Request `mapstructure:"requests" json:"requests" yaml:"requests"`
}

type Request struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}
