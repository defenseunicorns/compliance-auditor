# File Domain
The File domain allows for validation of arbitrary file contents. The file domain can evaluate local files and network files. Files are copied to a temporary directory for evaluation and deleted afterwards.

## Specification
The File domain specification accepts a descriptive name for the file as well as it's path:

```yaml
domain:
  type: file
  file-spec:
    filepaths:
    - name: config
      path: grafana.ini
```

## Supported File Types
The file domain use's OPA's [conftest](https://conftest.dev) to parse files into a json-compatible format for validations.  Both OPA and kyverno (using [kyverno-json](https://kyverno.github.io/kyverno-json/latest/)) can validate files parsed by the file domain.

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