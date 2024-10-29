# API Domain

The API Domain allows for collection of data (via HTTP Get Requests) generically from API endpoints. 

# Specification
The API domain Specification accepts a list of `Requests` and an `Options` block. `Options` can be configured at the top-level and will apply to all requests except those which have embedded `Options`. `Request`-level options will *override* top-level `Options`.


```yaml
domain: 
    type: api
    api-spec:
        # Options specified at this level will apply to all requests except those with an Options block of their own.
        Options:
            # Timeout configures the Request Timeout. Default is no timeout. The Timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
            Timeout: 30s
            # Configure a Proxy server for all requests.
            Proxy: "https://my.proxy"
            Headers: 
                key: "value"
                my-customer-header: "my-custom-value"
        #  Requests is a list of requests. The Request Name is the key used when referencing the resources returned from the API in the provider.
        Requests:
            # User-defined descriptive name
            - Name: "healthcheck" 
            # The URL of the Request. The API domain supports be any rfc3986-formatted URI. Lula also supports url `Parameters` as a separate argument. 
              Url: "https://example.com/health/ready"
            # Url Parameters to append to the URL. Lula also supports full URIs in the `Url`.
              Parameters: 
                key: "value"
              # Request-level Options have the same specification as the api-spec-level Options. These options apply only to this request.
              Options:
                # Timeout configures the Request Timeout. Default is no timeout. The Timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
                Timeout: 30s
                # Configure a Proxy server for all requests.
                Proxy: "https://my.proxy"
                Headers: 
                    key: "value"
                    my-customer-header: "my-custom-value"
            - Name: "readycheck"
            # etc ...
```
