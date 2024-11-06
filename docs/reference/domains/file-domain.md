# File Domain
The File domain allows for validation of arbitrary file contents from a list of supported file types. The file domain can evaluate local files and network files. Files are copied to a temporary directory for evaluation and deleted afterwards.

## Specification
The File domain specification accepts a descriptive name for the file as well as its path. The filenames and descriptive names must be unique.

```yaml
domain:
  type: file
  file-spec:
    filepaths:
    - name: config
      path: grafana.ini
      parser: ini         # optionally specify which parser to use for the file type
```

## Supported File Types
The file domain uses OPA's [conftest](https://conftest.dev) to parse files into a json-compatible format for validations. Both OPA and Kyverno (using [kyverno-json](https://kyverno.github.io/kyverno-json/latest/)) can validate files parsed by the file domain.

The file domain includes the following file parsers:
* cue
* cyclonedx
* dockerfile
* dotenv
* edn
* hcl1
* hcl2
* hocon
* ignore
* ini
* json
* jsonc
* jsonnet
* properties
* spdx
* string
* textproto
* toml
* vcl
* xml
* yaml

The file domain can also parse arbitrary file types as strings. The entire file contents will be represented as a single string.

The file parser can usually be inferred from the file extension. However, if the file extension does not match the filetype you are parsing (for example, if you have a json file that does not have a `.json` extension), or if you wish to parse an arbitrary file type as a string, use the `parser` field in the FileSpec to specify which parser to use. The list above contains all the available parses. 

## Validations
When writing validations against files, the filepath `name` must be included as
the top-level key in the validation. The placement varies between providers.

Given the following ini file:

```grafana.ini
[server]
# Protocol (http, https, socket)
protocol = http
```

The below Kyverno policy validates the protocol is https by including Grafana as the top-level key under `check`:

```yaml
metadata:
  name: check-grafana-protocol
  uuid: ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
domain:
  type: file
  file-spec:
    filepaths:
    - name: 'grafana'
      path: 'custom.ini'
provider:
  type: kyverno
  kyverno-spec:
    policy:
      apiVersion: json.kyverno.io/v1alpha1
      kind: ValidatingPolicy
      metadata:
        name: grafana-config
      spec:
        rules:
        - name: protocol-is-https
          assert:
            all:
            - check:
                grafana:
                  server:
                    protocol: https
```

While in an OPA policy, the filepath `Name` is the input key to access the config:

```yaml
metadata:
  name: validate-grafana-config
  uuid: ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
domain:
  type: file
  file-spec:
    filepaths:
    - name: 'grafana'
      path: 'custom.ini'
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      # Default values
      default validate := false
      default msg := "Not evaluated"
      
      validate if {
       check_grafana_config.result
      }
      msg = check_grafana_config.msg

      config := input["grafana"]
      protocol := config.server.protocol

      check_grafana_config = {"result": true, "msg": msg} if {
        protocol == "https"
        msg := "Server protocol is set to https"
      } else = {"result": false, "msg": msg} if {
        protocol == "http"
        msg := "Grafana protocol is insecure"
      }

    output:
      validation: validate.validate
      observations:
        - validate.msg
```

### Parsing files as arbitrary strings
Files that are parsed as strings are represented as a key-value pair where the key is the user-supplied file `name` and the value is a string representation of the file contexts, including special characters, for e.g. newlines (`\n`). 

As an example, let's parse a similar file as before as an arbitrary string. 

When reading the following multiline file contents as a string:
```server.txt
server = https
port = 3000
```

The resources for validation will be formatted as a single string with newline characters:

```
{"config": "server = https\nport = 3000"}
```

And the following validation will confirm if the server is configured for https:
```validation.yaml
  domain:
    type: file
    file-spec:
      filepaths:
      - name: 'config'
        path: 'server.txt'
        parser: string
  provider:
    type: opa
    opa-spec:
      rego: |
        package validate
        import rego.v1

        # Default values
        default validate := false
        default msg := "Not evaluated"
        
        validate if {
          check_server_protocol.result
        }
        msg = check_server_protocol.msg
        
        config := input["config"]
        
        check_server_protocol = {"result": true, "msg": msg} if {
          regex.match(
            `server = https\n`,
            config
          )
          msg := "Server protocol is set to https"
        } else = {"result": false, "msg": msg} if {
          regex.match(
            `server = http\n`,
            config
          )
          msg := "Server Protocol must be https - http is disallowed"
        }

      output:
        validation: validate.validate
        observations:
          - validate.msg
```

## Note on Compose
While the file domain is capable of referencing relative file paths in the `file-spec`, Lula does not de-reference those paths during composition. If you are composing multiple files together, you must either use absolute filepaths (including network filepaths), or ensure that all referenced filepaths are relative to the output directory of the compose command. 

## Evidence Collection

The use of `lula dev get-resources` and `lula validate --save-resources` will produce evidence in the form of `json` files. These files provide point-in-time evidence for auditing and review purposes.

Evidence collection occurs for each file specified and - in association with any error will produce an empty representation of the target file(s) data to be collected.
