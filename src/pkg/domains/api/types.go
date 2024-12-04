package api

import (
	"context"
	"net/url"
	"time"

	"github.com/defenseunicorns/lula/src/types"
)

// ApiDomain is a domain that is defined by a list of API requests
type ApiDomain struct {
	// Spec is the specification of the API requests
	Spec *ApiSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func CreateApiDomain(spec *ApiSpec) (types.Domain, error) {
	// Check validity of spec
	err := validateAndMutateSpec(spec)
	if err != nil {
		return nil, err
	}

	return ApiDomain{
		Spec: spec,
	}, nil
}

func (a ApiDomain) GetResources(ctx context.Context) (types.DomainResources, error) {
	return a.makeRequests(ctx)
}

// IsExecutable returns true if any of the requests are marked executable
func (a ApiDomain) IsExecutable() bool {
	return a.Spec.executable
}

// ApiSpec contains a list of API requests
type ApiSpec struct {
	Requests []Request `mapstructure:"requests" json:"requests" yaml:"requests"`
	// Opts will be applied to all requests, except those which have their own
	// specified ApiOpts
	Options *ApiOpts `mapstructure:"options" json:"options,omitempty" yaml:"options,omitempty"`

	// internally-managed fields executable will be set to true during spec
	// validation if *any* of the requests are flagged executable
	executable bool
}

// Request is a single API request
type Request struct {
	Name       string            `json:"name" yaml:"name"`
	URL        string            `json:"url" yaml:"url"`
	Params     map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Method     string            `json:"method,omitempty" yaml:"method,omitempty"`
	Body       string            `json:"body,omitempty" yaml:"body,omitempty"`
	Executable bool              `json:"executable,omitempty" yaml:"executable,omitempty"`
	// ApiOpts specific to this request. If ApiOpts is present, values in the
	// ApiSpec-level Options are ignored for this request.
	Options *ApiOpts `json:"options,omitempty" yaml:"options,omitempty"`

	// internally-managed options
	reqURL        *url.URL
	reqParameters url.Values
}

type ApiOpts struct {
	// Timeout in seconds
	Timeout string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Proxy   string            `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`

	// internally-managed options
	timeout  *time.Duration
	proxyURL *url.URL
}
