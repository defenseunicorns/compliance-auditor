package api

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

var defaultTimeout = 30 * time.Second

// validateAndMutateSpec validates the spec values and applies any defaults or
// other mutations or normalizations necessary. The original values are not modified.
// validateAndMutateSpec will validate the entire object and may return multiple
// errors.
func validateAndMutateSpec(spec *ApiSpec) (errs error) {
	if spec == nil {
		return errors.New("spec is required")
	}
	if len(spec.Requests) == 0 {
		errs = errors.Join(errs, errors.New("some requests must be specified"))
	}

	if spec.Options == nil {
		spec.Options = &ApiOpts{}
	}
	err := validateAndMutateOptions(spec.Options)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	for i := range spec.Requests {
		if spec.Requests[i].Name == "" {
			errs = errors.Join(errs, errors.New("request name cannot be empty"))
		}
		if spec.Requests[i].URL == "" {
			errs = errors.Join(errs, errors.New("request url cannot be empty"))
		}
		reqUrl, err := url.Parse(spec.Requests[i].URL)
		if err != nil {
			errs = errors.Join(errs, errors.New("invalid request url"))
		} else {
			spec.Requests[i].reqURL = reqUrl
		}
		if spec.Requests[i].Params != nil {
			queryParameters := url.Values{}
			for k, v := range spec.Requests[i].Params {
				queryParameters.Add(k, v)
			}
			spec.Requests[i].reqParameters = queryParameters
		}
		if spec.Requests[i].Options != nil {
			err = validateAndMutateOptions(spec.Requests[i].Options)
			if err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}

	return errs
}

func validateAndMutateOptions(opts *ApiOpts) (errs error) {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}

	if opts.Timeout != "" {
		duration, err := time.ParseDuration(opts.Timeout)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("invalid wait timeout string: %s", opts.Timeout))
		}
		opts.timeout = &duration
	}

	if opts.timeout == nil {
		opts.timeout = &defaultTimeout
	}

	if opts.Proxy != "" {
		proxyURL, err := url.Parse(opts.Proxy)
		if err != nil {
			// not logging the input URL in case it has embedded credentials
			errs = errors.Join(errs, errors.New("invalid proxy string"))
		}
		opts.proxyURL = proxyURL
	}

	return errs
}
