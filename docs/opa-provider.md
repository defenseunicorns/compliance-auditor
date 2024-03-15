# OPA Provider

The OPA provider provides Lula with the capability to evaluate the `domain` in target against a rego policy. 

## Payload Expectation

The validation performed should be in the form of provider, domain, and payload.

Example:
```yaml
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
    rego: |                                   # Required - Rego policy used for data validation
      package validate                        # Required - Package name

      import future.keywords.every            # Optional - Any imported keywords

      validate {                              # Required - Rule Name for evaluation - "validate" is the only supported rule
        every pod in input.podsvt {
          podLabel == "bar"
        }
      }
```

Optionally, an `output` can be specified in the `payload`. Currently, the default validation allowance/denial is given by `validate.validate`, which is really of the structure `<package-name>.<json-path-to-boolean-variable>`. If you have a desired alternative validation boolean variable, as well as additional observations to include, an output can be added to payload such as:
```yaml
target:
  provider: "opa"
  domain: "kubernetes"
  payload:
    resource-rules: 
    - group: 
      version: v1 
      resource: pods
      namespaces: [validation-test] 
    rego: |
      package mypackage 

      result {
        input.kind == "Pod"
        podLabel := input.metadata.labels.foo
        podLabel == "bar"
      }
      test := "my test string"
    output:
      validation: mypackage.result
      observations:
      - validate.test
```
The `validatation` field must specify a json path that resolves to a boolean value. The `observations` array currently only support variables that resolve as strings. These observations will be printed out in the `remarks` section of `relevant-evidence` in the assessment results.

## Policy Creation

The required structure for writing a validation in rego for Lula to validate is as follows:

```yaml
rego: |
  package validate 

  validate {

  }
```

This structure can be utilized to evaluate an expression directly:

```yaml
rego: |
  package validate 

  validate {
    input.kind == "Pod"
    podLabel := input.metadata.labels.foo
    podLabel == "bar"
  }
```

The expression can also use multiple `rule bodies` as such:

```yaml
rego: |
  package validate

  foolabel {
    input.kind == "Pod"
    podLabel := input.metadata.labels.foo
    podLabel == "bar"
  }

  validate {
    foolabel
  }
```

> [!IMPORTANT]
> `package validate` and `validate` are required package and rule for Lula use currently. 
