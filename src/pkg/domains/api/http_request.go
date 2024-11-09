package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func doHTTPReq[T any](ctx context.Context, client http.Client, url url.URL, headers map[string]string, queryParameters url.Values, respTy T) (T, int, error) {
	// append any query parameters.
	q := url.Query()

	for k, v := range queryParameters {
		// using Add instead of set incase the input URL already had a query encoded
		q.Add(k, strings.Join(v, ","))
	}
	// set the query to the encoded parameters
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return respTy, 0, err
	}
	// add each header to the request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// do the thing
	res, err := client.Do(req)
	if err != nil {
		return respTy, 0, err
	}
	if res == nil {
		return respTy, 0, fmt.Errorf("error: calling %s returned empty response", url.Redacted())
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return respTy, 0, err
	}

	var responseObject T
	err = json.Unmarshal(responseData, &responseObject)
	return responseObject, res.StatusCode, err
}

func clientFromOpts(opts *ApiOpts) http.Client {
	transport := &http.Transport{}
	if opts.proxyURL != nil {
		transport.Proxy = http.ProxyURL(opts.proxyURL)
	}
	c := http.Client{Transport: transport}
	if opts.timeout != nil {
		c.Timeout = *opts.timeout
	}
	return c
}
