package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/defenseunicorns/lula/src/types"
)

func (a ApiDomain) makeRequests(_ context.Context) (types.DomainResources, error) {
	collection := make(map[string]interface{}, 0)

	var defaultOpts *ApiOpts
	if a.Spec.Options == nil {
		defaultOpts = new(ApiOpts)
	} else {
		defaultOpts = a.Spec.Options
	}

	// configure the default HTTP client using any top-level Options. Individual requests with overrides will get bespoke clients.
	transport := &http.Transport{}
	if defaultOpts.Proxy != "" {
		proxy, err := url.Parse(a.Spec.Options.Proxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy url: %s", err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	defaultClient := &http.Client{Transport: transport}
	if defaultOpts.Timeout != 0 {
		defaultClient.Timeout = time.Duration(defaultOpts.Timeout) * time.Second
	}

	for _, request := range a.Spec.Requests {
		resp, err := defaultClient.Get(request.URL)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil,
				fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "application/json" {

			var prettyBuff bytes.Buffer
			err := json.Indent(&prettyBuff, body, "", "  ")
			if err != nil {
				return nil, err
			}
			prettyJson := prettyBuff.String()

			var tempData interface{}
			err = json.Unmarshal([]byte(prettyJson), &tempData)
			if err != nil {
				return nil, err
			}
			collection[request.Name] = tempData

		} else {
			return nil, fmt.Errorf("content type %s is not supported", contentType)
		}
	}
	return collection, nil
}
