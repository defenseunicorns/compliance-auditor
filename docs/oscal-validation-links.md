# Validation Identifiers

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
