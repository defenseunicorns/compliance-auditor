package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

func doHTTPReq[T any](ctx context.Context, client http.Client, method string, url url.URL, body io.Reader, headers map[string]string, queryParameters url.Values, respTy T) (T, int, error) {
	if method == HTTPMethodGet {
		// append any query parameters
		q := url.Query()
		for k, v := range queryParameters {
			// using Add instead of set in case the input URL already had a query encoded
			q.Add(k, strings.Join(v, ","))
		}
		// set the query to the encoded parameters
		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return respTy, 0, err
	}
	// add each header to the request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// log the request
	message.Debug("%q %s", method, req.URL.Redacted())

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
