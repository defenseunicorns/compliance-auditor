package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/defenseunicorns/lula/src/types"
)

func (a ApiDomain) makeRequests(ctx context.Context) (types.DomainResources, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled: %s", ctx.Err())
	default:
		collection := make(map[string]interface{}, 0)

		// defaultOpts apply to all requests, but may be overridden by adding an
		// options block to an individual request.
		var defaultOpts *ApiOpts
		if a.Spec.Options == nil {
			// This isn't likely to be nil in real usage, since CreateApiDomain
			// parses and mutates specs.
			defaultOpts = new(ApiOpts)
			defaultOpts.timeout = &defaultTimeout
		} else {
			defaultOpts = a.Spec.Options
		}

		// configure the default HTTP client using any top-level Options. Individual
		// requests with overrides (in request.Options.Headers) will get bespoke clients.
		defaultClient := clientFromOpts(defaultOpts)
		var errs error
		for _, request := range a.Spec.Requests {
			var responseType map[string]interface{}

			var r io.Reader
			if request.Body != "" {
				r = bytes.NewBufferString(request.Body)
			}

			var headers map[string]string
			var client http.Client

			if request.Options == nil {
				headers = defaultOpts.Headers
				client = defaultClient
			} else {
				headers = request.Options.Headers
				client = clientFromOpts(request.Options)
			}

			responseType, status, err := doHTTPReq(ctx, client, request.Method, *request.reqURL, r, headers, request.reqParameters, responseType)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			// Check if the response object is empty and manually add a DR with the status response if so. This is more likely to happen in tests than reality.
			if responseType != nil {
				responseType["status"] = status
				collection[request.Name] = responseType
			} else {
				collection[request.Name] = types.DomainResources{"status": status}
			}
		}
		return collection, errs
	}
}
