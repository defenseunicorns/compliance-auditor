# Kyverno Provider

The Kyverno provider provides Lula with the capability to evaluate the `domain` in target against a rego policy. 

## Payload Expectation

The validation performed should be in the form of provider, domain, and payload.

Example:
```yaml
target:
  provider: kyverno
  domain: kubernetes
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

You can have mutiple policies defined. Optionally, `output.validation` can be specified in the `payload` to control which (Policy, Rule) pair control validation allowance/denial, which is in the structure of a comma separated list of rules: `policy-name1.rule-name-1,policy-name-1.rule-name-2`. If you have a desired observation to include, `output.observations` can be added to payload to observe violations by a certain (Policy, Rule) pair such as:
```yaml
target:
  provider: "kyverno"
  domain: "kubernetes"
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
    output:
      validation: labels.foo-label-exists
      observations:
      - labels.foo-label-exists
```
The `validatation` and `observations` fields must specify a (Policy, Rule) pair. These observations will be printed out in the `remarks` section of `relevant-evidence` in the assessment results.
