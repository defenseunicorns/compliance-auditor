# Validation Identifiers

- [Validation Identifiers](#validation-identifiers)
  - [Connecting Links with Lula Validations](#connecting-links-with-lula-validations)
    - [Rel](#rel)
  - [Importing Validations](#importing-validations)
    - [Local Validations](#local-validations)
    - [Remote Validations](#remote-validations)
    - [Checksums](#checksums)
    - [Multiple Validations](#multiple-validations)
___
In OSCAL - `links` contains the following fields:
```yaml
links:
  - href: https://www.example.com/
    rel: reference
    text: Example
    media-type: text/html
    resource-fragment: some-fragment
```

These links are a "reference to a local or remote resource, that has a specific relation to the containing object" - [Component Definition Links](https://pages.nist.gov/OSCAL-Reference/models/v1.1.2/component-definition/json-reference/#/component-definition/components/links).

As such, links are a native OSCAL attribute that Lula can use to map to Validations. 

## Connecting Links with Lula Validations

After identifying a control and writing a Lula Validation, we need to store that Lula Validation within the OSCAL artifact for referencing.

This is accomplished by adding a new `resource` to the `back-matter` as shown below:

```yaml
back-matter:
  resources:
  - uuid: a7377430-2328-4dc4-a9e2-b3f31dc1dff9
    description: >-
      domain:
        type: kubernetes
        kubernetes-spec:
          resources:
          - name: podsvt 
            resource-rule:   
              group: 
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

Now we need to map an existing control (or Component-Definition Implemented-Requirement) to this Lula Validation. 

### Rel
The default workflow is to use the rel attribute to indicate that Lula has work to perform.

In the instance of a standard validation - A link to a Lula Validation might look like this:
```yaml
links:
  - href: '#a7377430-2328-4dc4-a9e2-b3f31dc1dff9'
    rel: lula
```

Where `href: '#a7377430-2328-4dc4-a9e2-b3f31dc1dff9'` points to an OSCAL object with a UUID reference and `rel: lula` indicates that the link is to a Lula Validation.
UUID's should always be unique per object in the OSCAL artifact.


> [!TIP]
> You can generate a random UUID using `lula tools uuidgen` or a deterministic UUID using `lula tools uuidgen <string>`.

## Importing Validations
In addition to storing validaitons in the `BackMatter`, `links` may be used to fetch resources external to the `component-definition`.

### Local Validations
- must be prefixed with `file:`
- `file:` must be a relative path to the `component-definition` or an absolute path
```yaml
links:
  - href: file:./validation.yaml
    rel: lula
  - href: file:/home/user/validations/validation.yaml
    rel: lula
```

### Remote Validations
- must be prefixed with `https:` or `http:`
- `https:` or `http:` must be a valid URL
```yaml
links:
  - href: https://example.com/validation.yaml
    rel: lula
```

### Checksums
- A checksum may be provided in the href using the suffix `@<checksum>` 
- Supports `sha1`, `sha256`, `sha512`, `md5`
```yaml
links:
  - href: https://example.com/validation.yaml@0123456789abcdef
    rel: lula
```

### Multiple Validations 
- A file with multiple validations may be provided in the link.
- `---` should be used to separate each validation
- `resource-fragment: <UUID>` will run the validation with the UUID specified
- `resource-fragment: *` will run all validations
```yaml
// Only runs the validation with the UUID of a7377430-2328-4dc4-a9e2-b3f31dc1dff9
links:
  - href: https://example.com/multi-validations.yaml
    rel: lula
    resource-fragment: '#a7377430-2328-4dc4-a9e2-b3f31dc1dff9'
// All validations
  - href: file:./multi-validations.yaml
    rel: lula
    resource-fragment: *
```
___ 
> [!NOTE]
> An example `component-definition` with remote validations can be found [here](../src/test/e2e/scenarios/remote-validations/component-definition.yaml).
