package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/defenseunicorns/lula/src/types"
)

type APIResponse struct {
	StatusCode int
	Status     string
	Raw        json.RawMessage
	Response   any
}

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

			response, err := doHTTPReq(ctx, client, request.Method, *request.reqURL, r, headers, request.reqParameters)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			if response != nil {
				dr := types.DomainResources{
					"status":     response.Status,
					"statuscode": response.StatusCode,
					"raw":        response.Raw,
					"response":   response.Response,
				}
				collection[request.Name] = dr
			} else {
				// If the entire response is empty, return a validly empty resource
				collection[request.Name] = types.DomainResources{"status": 0}
			}
		}
		return collection, errs
	}
}
