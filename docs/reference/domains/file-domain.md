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
```

## Supported File Types
The file domain uses OPA's [conftest](https://conftest.dev) to parse files into a json-compatible format for validations. ∑Both OPA and Kyverno (using [kyverno-json](https://kyverno.github.io/kyverno-json/latest/)) can validate files parsed by the file domain.

The file domain supports the following file formats for validation:
* CUE
* CycloneDX
* Dockerfile
* EDN
* Environment files (.env)
* HCL and HCL2
* HOCON
* Ignore files (.gitignore, .dockerignore)
* INI
* JSON
* Jsonnet
* Property files (.properties)
* SPDX
* TextProto (Protocol Buffers)
* TOML
* VCL
* XML
* YAML

## Validations
When writing validations against files, the filepath `Name` must be included as
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

## Note on Compose
While the file domain is capable of referencing relative file paths in the `file-spec`, Lula does not de-reference those paths during composition. If you are composing multiple files together, you must either use absolute filepaths (including network filepaths), or ensure that all referenced filepaths are relative to the output directory of the compose command. 
