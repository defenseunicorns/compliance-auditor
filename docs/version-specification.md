# Version Specification
In cases where a specific version of Lula is desired, either for typing constraints or desired functionality, a `lula-version` property is recognized in the `description` (component-definition.back-matter.resources[_]):
```yaml
- uuid: 88AB3470-B96B-4D7C-BC36-02BF9563C46C
  title: Lula Validation
  remarks: >-
    No outputs in payload
  description: |
    lula-version: ">=0.0.2"
    target:
      provider: opa
      domain: kubernetes
      payload:
        resources:
        - name: podsvt
          resource-rule:
            group:
            version: v1
            resource: pods
            namespaces: [validation-test]
        rego: |                                   
          package validate

          import future.keywords.every

          validate { 
            every pod in input.podsvt {
              podLabel == "bar"
            }
          }
```

If included, the `lula-version` must be a string and should indicate the version constraints desired, if any. Our implementation uses Hashicorp's [go-version](https://pkg.go.dev/github.com/hashicorp/go-version) library, and constraints should follow their [conventions](https://developer.hashicorp.com/terraform/language/expressions/version-constraints). 

If an invalid string is passed or the current Lula version does not meet version constraints, the implementation will automatically be marked "not-satisfied" and a remark will be created in the Assessment Report detailing the rationale.