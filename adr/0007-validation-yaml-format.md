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
- Validations that call to action human evaluators (perhaps an additional provider/domain?)
- Logical validation compositions - saying that validations compose in an "OR" as opposed to default "AND" (this might also be external, more of a generation thing? Could be determined by the metadata.type?)

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
The following yaml is the proposed high-level structure for the validation file. The x_spec field under domain and provider is intended to be optional but should be populated for the selected type. The rationale for having different specs all wrapped into a single definition is to make it easier to unmarshal the entire yaml document at once, as opposed to having to piecemeal it. The spec that gets used in the validation logic is that which corresponds to the type chosen.

```yaml
lula-version: "1.0"                           # Optional (maintains backward compatilibity)
metadata:                                     # Optional
  name: "title here"                          # Optional (short description to use in output of validations could be useful)
  type: satisfaction                          # Optional (enum:[satisfaction, healthcheck, ?]) - basically this indicates how the validation is reported in results, default is probably just satisfaction, but this could add some extensibility to having various workflows depending on "type" values
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
    type: kyverno                             # Required (enum:[opa, kyverno])
    kyverno-spec:                             # Optional
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

Example for passthrough or placeholder validation: 

```yaml
target:
  domain: 
    type: passthrough                         # Required (enum:[kubernetes, passthrough])
    passthrough-spec:
      evaluation: INTERVIEW                   # Required: INTERVIEW | EXAMINE
  provider:
    type: human
    human-spec:
      guidance: "some guidance for user"
```