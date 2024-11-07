package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func doHTTPReq[T any](client http.Client, url url.URL, headers map[string]string, queryParameters url.Values, respTy T) (T, error) {
	// append any query parameters.
	q := url.Query()

	for k, v := range queryParameters {
		// using Add instead of set incase the input URL already had a query encoded
		q.Add(k, strings.Join(v, ","))
	}
	// set the query to the encoded parameters
	url.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return respTy, err
	}

	// add each header to the request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// do the thing
	res, err := client.Do(req)
	if err != nil {
		return respTy, err
	}

	if res == nil {
		return respTy, fmt.Errorf("error: calling %s returned empty response", url.Redacted())
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return respTy, err
	}

	if res.StatusCode != http.StatusOK {
		return respTy, fmt.Errorf("expected status code 200 but got %d", res.StatusCode)
	}

	var responseObject T
	err = json.Unmarshal(responseData, &responseObject)

	if err != nil {
		return respTy, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return responseObject, nil
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
