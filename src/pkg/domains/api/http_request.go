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

func doHTTPReq(ctx context.Context, client http.Client, method string, url url.URL, body io.Reader, headers map[string]string, queryParameters url.Values) (*APIResponse, error) {
	// append any query parameters
	q := url.Query()
	for k, v := range queryParameters {
		// using Add instead of set in case the input URL already had a query encoded
		q.Add(k, strings.Join(v, ","))
	}
	// set the query to the encoded parameters
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		message.Debugf("error from http.NewRequestWithContext: %s", err)
		return nil, err
	}
	// add each header to the request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// log the request
	message.Debugf("%q %s", method, req.URL.Redacted())

	// do the thing
	res, err := client.Do(req)
	if err != nil {
		message.Debugf("error from client.Do: %s", err)
		return nil, err
	}
	if res == nil {
		message.Debug("empty response")
		return nil, fmt.Errorf("error: %s returned empty response", url.Redacted())
	}
	defer res.Body.Close()
	var respObj APIResponse
	respObj.StatusCode = res.StatusCode
	if res.Status == "" {
		respObj.Status = http.StatusText(res.StatusCode)
	} else {
		respObj.Status = res.Status
	}
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		message.Debugf("error reading response body: %s", err)
		return nil, err
	}

	if respObj.StatusCode >= http.StatusOK && respObj.StatusCode < http.StatusMultiStatus {
		respObj.Raw = responseData
		err = json.Unmarshal(responseData, &respObj.Response)
		if err != nil {
			message.Debugf("error unmarshalling response: %s", err)
			return nil, err
		}
	}
	return &respObj, nil
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
