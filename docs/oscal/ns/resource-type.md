# Resource-Type

The Back-Matter OSCAL structures are used to store a collection of resources that may be referenced from within the OSCAL document instance. For Lula's purposes, the `resources` stored in the back-matter are used in the following ways:

- To store the "Lula Validation" artifacts
- To store templated "Lula Validation" artifacts

To identify these types of resources, the `resource-type` prop is used.

## Example

For a valid Lula Validation, the `resource-type` prop would be set to `validation`:
```yaml
props:
  - name: resource-type
    ns: https://docs.lula.dev/oscal/ns
    value: "validation"
```

For a templated Lula Validation, the `resource-type` prop would be set to `validation-template`:
```yaml
props:
  - name: resource-type
    ns: https://docs.lula.dev/oscal/ns
    value: "validation-template"
```
