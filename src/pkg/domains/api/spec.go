package api

import (
	"errors"
	"fmt"
	"net/http"
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

	for _, request := range spec.Requests {
		if request.Name == "" {
			errs = errors.Join(errs, errors.New("request name cannot be empty"))
		}
		if request.URL == "" {
			errs = errors.Join(errs, errors.New("request url cannot be empty"))
		}
		url, err := url.Parse(request.URL)
		if err != nil {
			errs = errors.Join(errs, errors.New("invalid request url"))
		} else {
			request.reqURL = url
		}
		if request.Method != "" {
			switch m := request.Method; m {
			case "Get", "get", "GET":
				request.method = http.MethodGet
			case "Head", "head", "HEAD":
				request.method = http.MethodHead
			case "Post", "post", "POST":
				request.method = http.MethodPost
			case "Postform", "postform", "POSTFORM": // PostForm is not a separate HTTP method, but it uses a different function (http.PostForm)
				request.method = "POSTFORM"
			default:
				errs = errors.Join(errs, fmt.Errorf("unsupported method: %s", request.Method))
			}
		}
		if request.Options != nil {
			validateAndMutateOptions(request.Options)
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
