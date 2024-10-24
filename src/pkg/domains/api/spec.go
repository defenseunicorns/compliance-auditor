package api

import (
	"errors"
	"fmt"
	"time"
)

// validateAndMutateSpec validates the spec values and applies any defaults or
// other mutations necessary. The original values are not modified.
// validateAndMutateSpec will validate the entire object and may return multiple
// errors.
func validateAndMutateSpec(spec *ApiSpec) (errs error) {
	if spec == nil {
		return errors.New("spec is required")
	}
	if len(spec.Requests) == 0 {
		errs = errors.Join(errors.New("some requests must be specified"))
	}

	if spec.Options != nil {
		if spec.Options.Timeout != "" {
			duration, err := time.ParseDuration(spec.Options.Timeout)
			if err != nil {
				errs = errors.Join(fmt.Errorf("invalid wait timeout string: %s", spec.Options.Timeout))
			}
			spec.Options.timeout = &duration
		}
	} else {
		// add an Options struct with a sane default timeout
		spec.Options = &ApiOpts{}
	}
	if spec.Options.timeout == nil {
		d := 30 * time.Second
		spec.Options.timeout = &d
	}

	for _, request := range spec.Requests {
		if request.Name == "" {
			errs = errors.Join(errors.New("request name cannot be empty"))
		}
		if request.URL == "" {
			errs = errors.Join(errors.New("request url cannot be empty"))
		}
		if request.Method != "" {
			if request.Method != "Get" && request.Method != "Head" && request.Method != "Post" && request.Method != "PostForm" {
				errs = errors.Join(fmt.Errorf("unsupported method: %s", request.Method))
			}
		}
	}

	return errs
}
