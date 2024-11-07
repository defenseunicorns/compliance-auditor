# API Domain

The API Domain allows for collection of data (via HTTP Get Requests) generically from API endpoints. 

## Specification
The API domain Specification accepts a list of `Requests` and an `Options` block. `Options` can be configured at the top-level and will apply to all requests except those which have embedded `Options`. `Request`-level options will *override* top-level `Options`.


```yaml
domain: 
  type: api
  api-spec:
    # Options specified at this level will apply to all requests except those with an embedded options block.
    options:
      # Timeout configures the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
      timeout: 30s
      # Proxy specifies a proxy server for all requests.
      proxy: "https://my.proxy"
      # Headers is a map of key value pairs to send with all requests.
      headers: 
        key: "value"
        my-customer-header: "my-custom-value"
    # Requests is a list of URLs to query. The request name is the map key used when referencing the resources returned by the API.
    requests:
      # A descriptive name for the request.
      - name: "healthcheck" 
      # The URL of the request. The API domain supports any rfc3986-formatted URI. Lula also supports URL parameters as a separate argument. 
        url: "https://example.com/health/ready"
        # Parameters to append to the URL. Lula also supports full URIs in the URL.
        parameters: 
          key: "value"
        # Request-level options have the same specification as the api-spec-level options. These options apply only to this request.
        options:
          # Configure the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
          timeout: 30s
          # Proxy specifies a proxy server for this request.
          proxy: "https://my.proxy"
          # Headers is a map of key value pairs to send with this request.
          headers: 
            key: "value"
            my-customer-header: "my-custom-value"
      - name: "readycheck"
      # etc ...
```

## API Domain Resources

The API response body is serialized into a json object with the Request's Name as the top-level key. The API status code is included in the output domain resources.

Example output:

```json
"healthcheck": {
  "status": 200,
  "healthy": "ok"
}
```