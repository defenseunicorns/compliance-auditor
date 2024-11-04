package api

import (
	"context"
	"fmt"

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
		// requests with overrides will get bespoke clients.
		defaultClient := clientFromOpts(defaultOpts)

		for _, request := range a.Spec.Requests {
			var responseType interface{}
			var err error
			if request.Options == nil {
				responseType, err = doHTTPReq(ctx, defaultClient, *request.reqURL, defaultOpts.Headers, request.reqParameters, responseType)
			} else {
				client := clientFromOpts(request.Options)
				responseType, err = doHTTPReq(ctx, client, *request.reqURL, request.Options.Headers, request.reqParameters, responseType)
			}
			if err != nil {
				return collection, err
			}
			collection[request.Name] = responseType
		}
		return collection, nil
	}
}
