# 7. Validation YAML document schema

Date: 2024-04-02

## Status

Discussion

## Context

The format for a Lula Validation object (i.e., the Lula Validation yaml document) is possibly not extensible as Lula continues to grow, and also needs to track to the code structure to group functionality where it makes the most sense (e.g., resource collection). With this context in mind, the following are some guiding principles that should help us in developing this schema:
- Clarity and simplicity - the schema should be clear and intuitive to new users, we should use clear names to denote fields and avoid deep nestings
- Organization of the yaml should match with Lula architecture to both help Lula use the document and for navigability/readibility
- The schema should have versioning in place that supports changes and backward compatibility

The use cases for the validation yaml should support
- Local and remote validation files (this may be external to the validation yaml document itself)
- Different domains and providers, as well as their different specifications for use
- Validations that call to action human evaluators
- Logical validation compositions - saying that validations compose as "OR" as opposed to default "AND" 
- Validation overrides

### Current

The current yaml document options are as follows, depending on opa or kyverno providers

```yaml
lula-version: "1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [Kubernetes])
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
lula-version: "1.0"                           # Optional
target:
  provider: opa                               # Required (enum: [opa, kyverno])
  domain: kubernetes                          # Required (enum: [Kubernetes])
  payload:
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
The following yaml is the proposed high-level structure for the validation file. The X-spec field under domain and provider should be populated for the selected `type`. 
The rationale for having different specs in this format is to make it clear to the user which fields are relevant to the selected provider or domain. Previous impementation had them all more or less at the same level, so it might be confusing for a user to know which fields related to which domain or provider. Additionally, this allows for reusable property names across providers or domains.

Unfortunately this proposed structure adds more nesting, but hopefully the nesting structure isn't confusing to the user and instead provides a more obvious delineation on domain/provider specifications.

```yaml
lula-version: "1.0"                           # Optional (maintains backward compatilibity)
metadata:                                     # Optional
  name: "title here"                          # Optional (short description to use in output of validations could be useful)
target:
  domain: 
    type: kubernetes                          # Required (enum:[kubernetes, passthrough])
    kubernetes-spec:                          # Optional
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
target:
  domain: 
    type: kubernetes                          # Required (enum:[kubernetes, passthrough])
    kubernetes-spec:                          # Optional
      resources:                                  
      - name: podsvt                          # Required 
        resource-rule:                        # Required
          version: v1                         # Required
          kind: pods                          # Required (formerly "resource" but "kind" seems to make more sense in a k8s context)
          namespaces: [validation-test]       # Optional (Required with "name")
  provider: 
    type: kyverno                             # Required (enum:[opa, kyverno])
    kyverno-spec:                             # Optional
      kyverno:
        apiVersion: json.kyverno.io/v1alpha1    # Required
        kind: ValidatingPolicy                  # Required
        metadata:
          name: pod-policy                      # Required
        spec:
          rules:
            - name: no-latest                   # Required
              # Match payloads corresponding to pods
              match:                            # Optional
                any:                            # Assertion Tree
                - apiVersion: v1
                  kind: Pod
              assert:                           # Required
                all:                            # Assertion Tree
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

### Consequences
The intent of these changes is to decouple the domain and providers.This should result in is a architecture that supports variable domains and providers, such that they can be easily swapped out for one another. It may not be possible to *fully* decouple the two as, for instance, a rego policy will always need to have some notion of the structure of the inputs, however the idea is to make more defined interfaces at least.