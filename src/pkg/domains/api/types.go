package api

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/lula/src/types"
)

// ApiDomain is a domain that is defined by a list of API requests
type ApiDomain struct {
	// Spec is the specification of the API requests
	Spec *ApiSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func CreateApiDomain(spec *ApiSpec) (types.Domain, error) {
	// Check validity of spec
	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	if len(spec.Requests) == 0 {
		return nil, fmt.Errorf("some requests must be specified")
	}
	for _, request := range spec.Requests {
		if request.Name == "" {
			return nil, fmt.Errorf("request name cannot be empty")
		}
		if request.URL == "" {
			return nil, fmt.Errorf("request url cannot be empty")
		}
	}

	return ApiDomain{
		Spec: spec,
	}, nil
}

func (a ApiDomain) GetResources(_ context.Context) (types.DomainResources, error) {
	return MakeRequests(a.Spec.Requests)
}

func (a ApiDomain) IsExecutable() bool {
	// Domain is not currently executable
	return false
}

// ApiSpec contains a list of API requests
type ApiSpec struct {
	Requests []Request `mapstructure:"requests" json:"requests" yaml:"requests"`
}

// Request is a single API request
type Request struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}
