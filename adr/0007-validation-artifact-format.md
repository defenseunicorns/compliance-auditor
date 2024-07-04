# 7. Validation Artifact Schema

Date: 2024-04-02

## Status

Accepted

## Context

The format for a Lula Validation artifact must be extensible as Lula continues to grow. It needs to track to the code structure to group functionality where it makes the most sense (e.g., resource collection). With this context in mind, the following are some guiding principles that should help us in developing this schema:
- Clarity and simplicity - the schema should be clear and intuitive to new users, we should use clear names to denote fields and avoid deep nestings
- Organization of the format should match with Lula architecture to both help Lula use the document and for navigability/readability
- The schema should have versioning in place that supports changes and backward compatibility

The use cases for the validation should support:
- Embedded, local and remote validation artifacts
- Different domains and providers, as well as their different specifications for use

### Current

The current yaml document options are as follows, depending on opa or kyverno providers

```yaml
lula-version: "v0.1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [kubernetes])
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

```yaml
lula-version: "v0.1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [kubernetes])
  payload:
    resources:
    - name: podsvt
      resource-rule:
        group:
        version: v1
        resource: pods
        namespaces: [validation-test]
    kyverno:
      apiVersion: json.kyverno.io/v1alpha1
      kind: ValidatingPolicy
      metadata:
        name: labels
      spec:
        rules:
        - name: foo-label-exists
          assert:
            all:
            - check:
                ~.podsvt:
                  metadata:
                    labels:
                      foo: bar
```

### Proposal

The following artifact is the proposed high-level structure for the validation. The X-spec field under domain and provider should be populated for the selected `type`. 
The rationale for having different specs in this format is to make it clear to the user which fields are relevant to the selected provider or domain. Previous implementation had them all more or less at the same level, so it might be confusing for a user to know which fields related to which domain or provider. Additionally, this allows for reusable property names across providers or domains.

```yaml
lula-version: ""                            # Optional (maintains backward compatibility)
metadata:                                   # Optional
  name: "title here"                        # Optional (short description to use in output of validations could be useful)
domain: 
  type: kubernetes                          # Required (enum:[kubernetes, api])
  kubernetes-spec:                          # Optional
    resources:                                  
    - name: podsvt                          # Required 
      resource-rule:                        # Required
        name:                               # Optional (Required with "field")
        group:                              # Optional (not all k8s resources have a group, the main ones are "")
        version: v1                         # Required
        resource: pods                      # Required
        namespaces: [validation-test]       # Optional (Required with "name")
        field:                              # Optional 
          jsonpath:                         # Required
          type:                             # Optional 
          base64:                           # Optional 
    wait:                                   # Optional 
      condition: Ready                      # Optional 
      kind: pod/test-pod-wait               # Optional 
      namespace: validation-test            # Optional 
      timeout: 30s                          # Optional 
provider: 
  type: opa                                 # Required (enum:[opa, kyverno])
  opa-spec:                                 # Optional
    rego: |                                 # Required 
      package validate

      validate := False
      test := "test string"
    output:                                 # Optional
      validation: validate.validate         # Optional
      observations:                         # Optional
      - validate.test                         
```

Example for kyverno:

```yaml
lula-version: ""                            # Optional (maintains backward compatibility)
metadata:                                   # Optional
  name: "title here"                        # Optional (short description to use in output of validations could be useful)
domain: 
  type: kubernetes                          # Required (enum:[kubernetes, api])
  kubernetes-spec:                          # Optional
    resources:                                  
    - name: podsvt                          # Required 
      resource-rule:                        # Required
        version: v1                         # Required
        resource: pods                      # Required 
        namespaces: [validation-test]       # Optional (Required with "name")
provider: 
  type: kyverno                             # Required (enum:[opa, kyverno])
  kyverno-spec:                             # Optional
    policy:                                 # Required
      apiVersion: json.kyverno.io/v1alpha1
      kind: ValidatingPolicy
      metadata:
        name: labels
      spec:
        rules:
        - name: foo-label-exists
          assert:
            all:
            - check:
                ~.podsvt:
                  metadata:
                    labels:
                      foo: bar
```
### Consequences

- Reorganization of the validation components into a hopefully more intuitive structure for clarity around domain/provider specifications.
- Decoupling domain and provider implementations from the core validation process.
- Provides a structure for future domain and providers to be implemented.