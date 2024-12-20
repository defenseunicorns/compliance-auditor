# Test a Validation

Writing tests for a Lula Validation should be a key part of the validation development process. The purpose of testing is to ensure that the Domain is returning expected data AND the Provider is correctly interpretting and "validating" that data. The testing framework is valuable to both document the tests the domain/provider passes (e.g., to aid in validation review), as well as setting up a repeatable test suite for the validations to be verified when the environment changes.

## About

This document will guide you through the process of writing tests for a Lula Validation. It will build on the [Develop a Validation](./develop-a-validation.md) guide, so it is recommended to read that first. Additional documentation on the testing framework can be found in the [testing](../reference/testing.md) reference.

## Writing Tests for a Lula Validation

### Pre-Requisites

* Lula installed
* Lula Validation from the [Develop a Validation](./develop-a-validation.md) guide

### Steps

>[!NOTE]
> Demo files can be found in the lula repository under `demo/validation-tests`

1. Assume we have the following Lula Validation:

```yaml
metadata:
  name: check-podinfo-health
  uuid: ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: podinfoDeployment
        resource-rule:
          name: my-release-podinfo
          namespaces: [podinfo]
          group: apps
          version: v1
          resource: deployments
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      # Default values
      default validate := false
      default msg := "Not evaluated"

      # Validation result
      validate if {
        check_podinfo_healthy.result
      }
      msg = check_podinfo_healthy.msg

      check_podinfo_healthy = {"result": true, "msg": msg} if {
        input.podinfoDeployment.spec.replicas > 0
        input.podinfoDeployment.status.availableReplicas == input.podinfoDeployment.status.replicas
        msg := "Number of replicas > 0 and all replicas are available."
      } else = {"result": false, "msg": msg} if {
        msg := "Podinfo not available."
      }
    output:
      validation: validate.validate
      observations:
        - validate.msg
```

We'd like to verify that our rego policy is going to correctly evaluate the `podinfoDeployment` resource if it should change.

2. We need to identify the types of changes we could expect to occur to the `podinfoDeployment` resource:

* If the resource is not found, we expect the policy to be `not-satisfied`
* If the resource is found, but the number of replicas is 0, we expect the policy to be `not-satisfied`
* If the resource is found, and the number of replicas is greater than 0, but the available replicas are not equal to the requested replicas, we expect the policy to be `not-satisfied`

3. Now that we've enumerated the possible outcomes, we can write our tests. We'll start with the first test, which is to verify that the policy is `not-satisfied` if the resource is not found. We know that if the `podinfoDeployment` is not found, the following JSON will result from the domain spec:

```json
{
  "podinfoDeployment": {}
}
```

To mimic this json structure in our test, we need to add the following to the `changes` section:

```yaml
- path: podinfoDeployment
  type: delete
- path: "."
  type: add
  value-map:
    podinfoDeployment: {}
```

These changes generate the above json structure by first removing the `podinfoDeployment` from the resources, and then adding it back with an empty map.

> [!NOTE]
> This is an interesting case that highlights the limitations of the change types - due to the way the underlying merge functionality works, it is not possible to update a map with empty keys. If a key exists, the only way to set it as empty is to first delete it, and then add it back with an empty map.

So we can add the following to the `tests` section to the `validation.yaml`:

```yaml
tests:
  - name: missing-podinfo-deployment
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment
        type: delete
      - path: "."
        type: add
        value-map:
          podinfoDeployment: {}
```

4. For the second test case, we want to verify that the policy is `not-satisfied` if the resource is found, but the number of replicas is 0. This mimics a scenario where the deployment is in the cluster, but there are no pods.

An abridged version of the json manifest we expect for this scenario is:

```json
{
  "podinfoDeployment": {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
      "name": "podinfo"
      // Rest of the metadata
    },
    "spec": {
      "replicas": 0
      // Rest of the spec
    }
  }
}
```

On first glance, we might be tempted to set the `podinfoDeployment.spec.replicas` to 0 using the following change:

```yaml
# invalid change for our resource!
- path: podinfoDeployment.status.replicas
  type: update
  value: "0"
```

However, this will NOT correctly generate the expected json structure since the `replicas` field is a number, and not a string. Instead, we need to use the `value-map` change type, which allows us to set the value of a field to any type of value, as follows:

```yaml
- path: podinfoDeployment.status
  type: update
  value-map:
    replicas: 0
```

Now the tests, containing both test cases, will become:

```yaml
tests:
  - name: missing-podinfo-deployment
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment
        type: delete
      - path: "."
        type: add
        value-map:
          podinfoDeployment: {}
  - name: zero-replicas
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment.status
        type: update
        value-map: 
            replicas: 0
```

5. Finally, the last test case checks the scenario where replicas are available, but the number of expected replicas is not equal to the number of available replicas.

This case yeilds a structure similar to the previous case, where the change is:
```yaml
- path: podinfoDeployment.status
  type: update
  value-map:
    availableReplicas: 0
```

6. We can bring this back together and compose our validation:

```yaml
metadata:
  name: check-podinfo-health
  uuid: ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
      - name: podinfoDeployment
        resource-rule:
          name: my-release-podinfo
          namespaces: [podinfo]
          group: apps
          version: v1
          resource: deployments
provider:
  type: opa
  opa-spec:
    rego: |
      package validate
      import rego.v1

      # Default values
      default validate := false
      default msg := "Not evaluated"

      # Validation result
      validate if {
        check_podinfo_healthy.result
      }
      msg = check_podinfo_healthy.msg

      check_podinfo_healthy = {"result": true, "msg": msg} if {
        input.podinfoDeployment.spec.replicas > 0
        input.podinfoDeployment.status.availableReplicas == input.podinfoDeployment.status.replicas
        msg := "Number of replicas > 0 and all replicas are available."
      } else = {"result": false, "msg": msg} if {
        msg := "Podinfo not available."
      }
    output:
      validation: validate.validate
      observations:
        - validate.msg
tests:
  - name: missing-podinfo-deployment
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment
        type: delete
      - path: "."
        type: add
        value-map:
          podinfoDeployment: {}
  - name: zero-replicas
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment.spec
        type: update
        value-map: 
            replicas: 0
  - name: not-equal-replicas
    expected-result: not-satisfied
    changes:
      - path: podinfoDeployment.status
        type: update
        value-map: 
            availableReplicas: 0
```

7. Now that we have our validation and appropriate tests, we can run `lula dev validate` from our `demo/test-validation` directory.

```sh
lula dev validate -f validation.yaml -r resources.json --run-tests --print-test-resources
```

And we should see the following output:

```sh
  ✔  Pass: missing-podinfo-deployment
  •  Result: not-satisfied
  •  --> validate.msg: Podinfo not available.
  •  Test Resources File Path: missing-podinfo-deployment.json
  ✔  Pass: zero-replicas
  •  Result: not-satisfied
  •  --> validate.msg: Podinfo not available.
  •  Test Resources File Path: zero-replicas.json
  ✔  Pass: not-equal-replicas                 
  •  Result: not-satisfied
  •  --> validate.msg: Podinfo not available.
  •  Test Resources File Path: not-equal-replicas.json 
```

> [!NOTE]
> The `--print-test-resources` flag is useful for debugging, as it will print the resources used for each test to the validation directory.