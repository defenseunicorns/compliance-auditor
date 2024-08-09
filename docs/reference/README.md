# Validation Reference

### Validation Struct

The `Validation` struct is a data structure used for ingesting validation data. It contains the following fields:

- `LulaVersion` (string): Optional field to maintain backward compatibility.
- `Metadata` (*Metadata): Optional metadata containing the name and UUID of the validation.
- `Provider` (*Provider): Required field specifying the provider and its corresponding specification.
- `Domain` (*Domain): Required field specifying the domain and its corresponding specification.

#### Metadata Struct

The `Metadata` struct contains the following fields:

- `Name` (string): Optional short description to use in the output of validations.
- `UUID` (string): Optional UUID of the validation.

#### Domain Struct

The `Domain` struct contains the following fields:

- `Type` (string): Required field specifying the type of domain (enum: `kubernetes`, `api`).
- `KubernetesSpec` (*KubernetesSpec): Optional specification for a Kubernetes domain, required if type is `kubernetes`.
- `ApiSpec` (*ApiSpec): Optional specification for an API domain, required if type is `api`.

#### Provider Struct

The `Provider` struct contains the following fields:

- `Type` (string): Required field specifying the type of provider (enum: `opa`, `kyverno`).
- `OpaSpec` (*OpaSpec): Optional specification for an OPA provider.
- `KyvernoSpec` (*KyvernoSpec): Optional specification for a Kyverno provider.

### Example YAML Document

The following is an example of a YAML document for a validation artifact:
```yaml
lula-version: ">=v0.2.0"
metadata:
  name: Validate pods with label foo=bar
  uuid: 123e4567-e89b-12d3-a456-426655440000
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: podsvt
        resource-rule:
          version: v1
          resource: pods
          namespaces: [validation-test]
provider:
  type: opa
  opa-spec:
    rego: |
      package validate

      import future.keywords.every

      validate {
        every pod in input.podsvt {
          podLabel := pod.metadata.labels.foo
          podLabel == "bar"
        }
      }
```
## Linting
Linting is done by Lula when a `Validation` object is converted to a `LulaValidation` for evaluation.

The `common.Validation.Lint` method is a convenience method to lint a `Validation` object. It performs the following step:

1. **Marshalling**: The method marshals the `Validation` object into a YAML byte array using the `common.Validation.MarshalYaml` function.
2. **Linting**: The method runs linting against the marshalled `Validation` object. This is done using the `schemas.Validate` function, which ensures that the YAML data conforms to the expected [schema](https://raw.githubusercontent.com/defenseunicorns/lula/main/src/pkg/common/schemas/validation.json).

___
The `schemas.Validate` function is responsible for validating the provided data against a specified JSON schema using [github.com/santhosh-tekuri/jsonschema/v6](https://github.com/santhosh-tekuri/jsonschema). The process involves the following steps:

1. **Coercion to JSON Map**: The provided data, which can be either an interface or a byte array, is coerced into a JSON map using the `model.CoerceToJsonMap` function.
2. **Schema Retrieval**: The function retrieves the JSON schema specified by the `schema` parameter using the `GetSchema` function.
3. **Schema Compilation**: The retrieved schema is compiled into a format that can be used for validation using the `jsonschema.CompileString` function.
4. **Validation**: The coerced JSON map is validated against the compiled schema. If the validation fails, the function extracts the specific errors and returns them as a formatted string.

## VS Code intellisense:
1. Ensure that the [YAML (Red Hat)](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) extension is installed.
2. Add the following to your settings.json:
```json
"yaml.schemas": {
    "${PATH_TO_LULA}/lula/src/pkg/common/schemas/validation.json": "*validation*.yaml"
},
```


> **Note:**
> - `${PATH_TO_LULA}` should be replaced with your path.
> - `*validation*.yaml` may be changed to match your project's validation file naming conventions.
> - can also be limited to project or workspace settings if desired 