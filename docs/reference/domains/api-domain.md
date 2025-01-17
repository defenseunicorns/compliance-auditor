# API Domain

The API Domain allows for collection of data (via HTTP Get or Post Requests) generically from API endpoints.

>[!Important]
>This domain supports both read and write operations (Get and Post operations), so use with care. If you configure validations with API calls that mutate resources, add the `executable` flag to the request so that Lula will ask for verification before making the API call. 

## Specification
The API domain Specification (`api-spec`) accepts a list of `requests` and an `options` block. `options` can be configured at the top-level and will apply to all requests except those which have embedded `options`. `request`-level `options` will *override* top-level `options`.


```yaml
domain: 
  type: api
  api-spec:
    # options (optional): Options specified at this level will apply to all requests except those with an embedded options block.
    options:
      # timeout (optional, default 30s): configures the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
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
        # method (optional, default get): The HTTP Method to use for the API call. "get" and "post" are supported. Default is "get".
        method: "get"
        # parameters (optional): parameters to append to the URL. Lula also supports full URIs in the URL.
        parameters: 
          key: "value"
        # Body (optional): a json-compatible string to pass into the request as the request body.
        body: |
stringjsondata
        # executable (optional, default false): Lula will request user verification before performing API actions if *any* API request is flagged "executable".
        executable: true
        # options (optional): Request-level options have the same specification as the api-spec-level options at the top. These options apply only to this request.
        options:
          # timeout (optional, default 30s): configures the request timeout. The default timeout is 30 seconds (30s). The timeout string is a number followed by a unit suffix (ms, s, m, h, d), such as 30s or 1m.
          timeout: 30s
          # proxy (optional): Specifies a proxy server for this request.
          proxy: "https://my.proxy"
          # headers (optional): a map of key value pairs to send with this request.
          headers: 
            key: "value"
            my-customer-header: "my-custom-value"
      - name: "readycheck"
      # etc ...
```

## API Domain Resources

The API response body is serialized into a json object with the `request` `name` as the top-level key. The API status code is included in the output domain resources under `status`. `raw` contains the entire API repsonse in an unmarshalled (`json.RawMessage`) format.

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

You can define a variable to simplify referencing the response: 
```
provider: 
  type: opa
  opa-spec:
    rego: |
      package validate

      resp := input.healthcheck.response

      validate {
        resp == true
      }
```

## Chaining API Calls
(here for the POC; will be integrated into the main document when real)
Additional fields in the Spec support chaining API requests:

- use `Outputs` & kyaml syntax to store values from one API call & reference in another:
```yaml
requests:
  - name: keycloak-token
    outputs:
    - name: token
      path: // kyaml syntax, relative to the output of this request 
    url: //etc
```

Now future requests can use any of the "tpl" fields to specify templates (using go syntax, with [[ ]] as delimiters to differentiate from our whole-file templatizing)

Yes all the names are terrible I'm VERY open to suggestions. In the long term / real implementation I'd like this all to just happen if outputs are detected instead of using additional (`tpl`) fields, but again: poc 

```yaml
requests:
  - name: keycloak-realm-check
    options:
      headerstpl:
        # syntax is requestname.outputname
        authorization: "Bearer [[.keycloak-token.token]]"
```