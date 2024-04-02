# 7. Validation YAML document schema

Date: 2024-04-02

## Status

Discussion

## Context

The format for a Lula Validation object (i.e., the Lula Validation yaml document) is possibly not extensible as Lula continues to grow, and also needs to track to the code structure to group functionality where it makes the most sense (e.g., resource collection). With this context in mind, the following are some guiding principles that should help us in developing this schema:
- Clarity and simplicity - the schema should be clear and intuitive to new users, we should use clear names to denote fields and avoid deep nestings
- Organization of the yaml should match with Lula architecture to both help Lula use the document and for navigability/readibility
- The schema should have versioning in place that supports changes and backward compatibility

### Current

The current yaml document options are as follows, depending on opa or kyverno providers

```yaml
lula-version: "1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [Kubernetes])
  payload:
    # This is a Variable Structure...
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
lula-version: "1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [Kubernetes])
  payload:
    # This is a Variable Structure...
    resources:
    - name: podsvt
      resource-rule:
        group:
        version: v1
        resource: pods
        namespaces: [validation-test]
    kyverno:
      apiVersion: json.kyverno.io/v1alpha1          # Required
      kind: ValidatingPolicy                        # Required
      metadata:
        name: pod-policy                            # Required
      spec:
        rules:
          - name: no-latest                         # Required
            # Match payloads corresponding to pods
            match:                                  # Optional
              any:                                  # Assertion Tree
              - apiVersion: v1
                kind: Pod
            assert:                                 # Required
              all:                                  # Assertion Tree
              - message: Pod `{{ metadata.name }}` uses an image with tag `latest`
                check:
                  ~.podsvt:
                    spec:
                      # Iterate over pod containers
                      # Note the `~.` modifier, it means we want to iterate over array elements in descendants
                      ~.containers:
                        image:
                          # Check that an image tag is present
                          (contains(@, ':')): true
                          # Check that the image tag is not `:latest`
                          (ends_with(@, ':latest')): false
```

### Proposal
The first yaml document attempts to create groups for domain and providers that can be extensible as new domains or providers are added, whereby the "spec" is the only thing that changes. The second yaml document is an attempt to look at doing "remote" validation, I'm not sure how we could consolidate the local and remote into one, so proposing simply using a different lula "type". The third is thinking about the logical composition of validation artifacts, it could enable some complex logic or even hierarchies.

```yaml
lula-version: "1.0"                           # Optional (maintains backward compatilibity)
metadata:                                     # Required
  title: "title here"                         # Optional (short description to use in output of validations could be useful)
  type: local                                 # Required (enum:[local, remote, admin]??)
  workflow: healthcheck                       # Required (enum: [satisfaction, healthcheck]) - basically this indicates some way that this validation is mapped to results
spec:
  domain: kubernetes                          # Required (enum:[kubernetes])
    domain-spec:                              # Required (the spec itself can be dependent on the domain chosen, so can be flexible to new domains like infra)
      # Variable, per domain, the following is for kubernetes
      resources:                                  
      - name: podsvt                          # Required 
        resource-rule:                        # Required
          name:                               # Optional (Required with "field")
          group:                              # Optional (not all k8s resources have a group, the main ones are "")
          version: v1                         # Required
          kind: pods                          # Required (formerly "resource" but "kind" seems to make more sense in a k8s context)
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

  provider: opa                               # Required (enum:[opa, kyverno])
    provider-spec:                            # Required (again, the spec itself can be dependent on the provider chosen)
      # Variable, per provider, the following is for opa
      rego: |                                 # Required 
        package validate

        validate := False
        test := "test string"
      output:                                 # Optional
        validation: validate.validate         # Optional
        observations:                         # Optional
        - validate.test                         
```

```yaml
lula-version: "1.0"                           # Optional 
metadata:                                     # Required
  title: "title here"                         # Optional 
  type: remote                                # Required 
  workflow: satisfaction                      # Required 
spec:
  validation:
    path: https://github.com/defenseunicorns/lula-compliance-lib/blob/main/validation.yaml                 
```

```yaml
lula-version: "1.0"                           # Optional 
metadata:                                     # Required
  title: "title here"                         # Optional 
  type: admin                                 # Required 
  workflow: satisfaction                      # Required 
spec:
  logic:
    or:
      - validation:
          path: https://github.com/defenseunicorns/lula-compliance-lib/blob/main/validation1.yaml
      - validation:
          path: https://github.com/defenseunicorns/lula-compliance-lib/blob/main/validation2.yaml                  
```