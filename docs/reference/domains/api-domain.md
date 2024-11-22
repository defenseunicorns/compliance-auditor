# API Domain

The API Domain allows for collection of data (via HTTP Get or Post Requests) generically from API endpoints.

## Specification
The API domain Specification accepts a list of `Requests` and an `Options` block. `Options` can be configured at the top-level and will apply to all requests except those which have embedded `Options`. `Request`-level options will *override* top-level `Options`.


```yaml
domain: 
  type: api
  api-spec:
    # options (optional): Options specified at this level will apply to all requests except those with an embedded options block.
    options:
      # timeout (optional): configures the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
      timeout: 30s
      # proxy (optional): Specifies a proxy server for all requests.
      proxy: "https://my.proxy"
      # headers (optional):  a map of key value pairs to send with all requests.
      headers: 
        key: "value"
        my-customer-header: "my-custom-value"
    # Requests is a list of URLs to query. The request name is the map key used when referencing the resources returned by the API.
    requests:
      # name (required): A descriptive name for the request.
      - name: "healthcheck" 
        # url (required): The URL for the request. The API domain supports any rfc3986-formatted URI. Lula also supports URL parameters as a separate argument.
        url: "https://example.com/health/ready"
        # method (optional): The HTTP Method to use for the API call. "get" and "post" are supported. Default is "get".
        method: "get"
        # parameters (optional): parameters to append to the URL. Lula also supports full URIs in the URL.
        parameters: 
          key: "value"
        # Body (optional): a json-compatible string to pass into the request as the request body.
        body: |
stringjsondata
        # executable (optional): Lula will request user verification before performing API actions if *any* API request is flagged "executable".
        executable: true
        # options (optional): Request-level options have the same specification as the api-spec-level options at the top. These options apply only to this request.
        options:
          # timeout (optional): configures the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
          timeout: 30s
          # proxy (optional): Specifies a proxy server for all requests.
          proxy: "https://my.proxy"
          # headers (optional):  a map of key value pairs to send with all requests.
          headers: 
            key: "value"
            my-customer-header: "my-custom-value"
      - name: "readycheck"
      # etc ...
```

## API Domain Resources

The API response body is serialized into a json object with the Request's `Name` as the top-level key. The API status code is included in the output domain resources under `status`. `raw` contains the entire API repsonse in an unmarshalled (`json.RawMessage`) format.

Example output:

```json
"healthcheck": {
  "status": 200,
  "response": {
    "healthy": true,
  },
  "raw": {"healthy": true}
}
```

The following example validation verifies that the request named "healthcheck" returns `"healthy": true` 

```
provider: 
  type: opa
  opa-spec:
    rego: |
      package validate

      validate {
        input.healthcheck.response.healthy == true
      }
```