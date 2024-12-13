# OPA Provider

The OPA provider provides Lula with the capability to evaluate the `domain` against a rego policy. 

## Payload Expectation

The validation performed should use the form of provider with the `type` of `opa` and using the `opa-spec`, along with a valid domain.

Example:
```yaml
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
    rego: |                                   # Required - Rego policy used for data validation
      package validate                        # Required - Package name

      import future.keywords.every            # Optional - Any imported keywords

      validate {                              # Required - Rule Name for evaluation - "validate" is the only supported rule
        every pod in input.podsvt {
          podLabel == "bar"
        }
      }
```

Optionally, an `output` can be specified in the `opa-spec`. Currently, the default validation allowance/denial is given by `validate.validate`, which is really of the structure `<package-name>.<json-path-to-boolean-variable>`. If you have a desired alternative validation boolean variable, as well as additional observations to include, an output can be added such as:
```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resource-rules: 
    - group: 
      version: v1 
      resource: pods
      namespaces: [validation-test] 
provider:
  type: opa
  opa-spec:
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
> `package validate` and `validate` are required package and rule for Lula use currently when an output.validation value has not been set. 

## Reusing OPA modules

Custom OPA modules can be imported and referenced in the main rego module. The following example shows how to
import a custom module and use it in the main rego module:

Let's say we have a file that is called `lula.rego` with the following contents that verifies that pods have
the label `lula: "true"`:

```yaml
package lula.labels

import rego.v1

has_lula_label(pod) if {
    pod.metadata.labels.lula == "true"
}
```

We can import this module and use it in the main rego module as follows:

```yaml
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
    modules:
      lula.labels: lula.rego
    rego: |                                   # Required - Rego policy used for data validation
      package validate                        # Required - Package name

      import future.keywords.every            # Optional - Any imported keywords
      import data.lula.labels as lula_labels  # Optional - Import the custom module

      validate {                              # Required - Rule Name for evaluation - "validate" is the only supported rule
        every pod in input.podsvt {
          lula_labels.has_lula_label(pod)     # Use the rules defined in the custom module
        }
      }
```
> [!Note]
> The `validate.rego` module name is reserved for the main rego policy and cannot be used as a custom module name.
