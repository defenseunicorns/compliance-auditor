# File Domain
The File domain allows for validation of arbitrary file contents. The file domain can evaluate local files and network files. Files are copied to a temporary directory for evaluation and deleted afterwards.

## Specification
The File domain specification accepts a descriptive name for the file as well as it's path. The names must be unique.

```yaml
domain:
  type: file
  file-spec:
    filepaths:
    - name: config
      path: grafana.ini
```

## Supported File Types
The file domain use's OPA's [conftest](https://conftest.dev) to parse files into a json-compatible format for validations. ∑Both OPA and kyverno (using [kyverno-json](https://kyverno.github.io/kyverno-json/latest/)) can validate files parsed by the file domain.

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
When writing validations against files, the filepath Name must be included as the top-level key in the validation, in this example below `check`:

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

```grafana.ini
[server]
# Protocol (http, https, socket)
protocol = http
```