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
      resourceRule:
        Group:
        Version: v1
        Resource: pods
        Namespaces: [validation-test]
    rego: |                                   # Required - Rego policy used for data validation
      package validate                        # Required - Package name

      import future.keywords.every            # Optional - Any imported keywords

      validate {                              # Required - Rule Name for evaluation - "validate" is the only supported rule
        every pod in input.podsvt {
          podLabel := pod.metadata.labels.foo
          podLabel == "bar"
        }
      }
```


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
