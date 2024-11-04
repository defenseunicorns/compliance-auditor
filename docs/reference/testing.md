# Testing

Testing is a key part of Lula Validation development. Since the results of the Lula Validations are determined by the policy set by the `provider`, those policies must be tested to ensure they are working as expected.

## Validation Testing

In the Lula Validation, a `tests` property is used to specify the each test that should be performed against the validation. Each test is a map of the following properties:

- `name`: The name of the test
- `changes`: An array of changes or transformations to be applied to the resources used in the test validation
- `expected-result`: The expected result of the test - satisfied or not-satisfied

A change is a map of the following properties:

- `path`: The path to the resource to be modified. The path syntax is described below.
- `type`: The type of operation to be performed on the resource
    - `update`: (default) updates the resource with the specified value
    - `delete`: deletes the field specified
    - `add`: adds the specified value
- `value`: The value to be used for the operation (string)
- `value-map`: The value to be used for the operation (map[string]interface{})

An example of a test added to a validation is:

```yaml
domain:
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: podsvt
      resource-rule:
        version: v1
        resource: pods
        namespaces: [validation-test]
provider:
  type: opa
  opa-spec:
    rego: |
      package validate

      import future.keywords.every

      validate {
        every pod in input.podsvt {
          podLabel := pod.metadata.labels.foo
          podLabel == "bar"
        }
      }
tests:
  - name: modify-pod-label-not-satisfied
    expected-result: not-satisfied
    changes:
      - path: podsvt.[metadata.namespace=validation-test].metadata.labels.foo
        type: update
        value: baz
  - name: delete-pod-label-not-satisfied
    expected-result: not-satisfied
    changes:
      - path: podsvt.[metadata.namespace=validation-test].metadata.labels.foo
        type: delete
```

There are two tests here:
* The first test will locate the first pod in the `validation-test` namespace and update the label `foo` to `baz`. Then a `validate` will be executed against the modified resources. The expected result of this is that the validation will fail, i.e., will be `not-satisfied`, which would result in a successful test.
* The second test will locate the first pod in the `validation-test` namespace and delete the label `foo`, then proceed to validate the modified resources and compare to the expected result.

### Path Syntax

This feature uses the kyaml library to inject data into the resources, so the path syntax is based on this library. 

The path should be a "." delimited string that specifies the keys along the path to the resource seeking to be modified. In addition to keys, a list item can be specified by using the “[some-key=value]” syntax. For example, the following path:

```
pods.[metadata.namespace=grafana].spec.containers.[name=istio-proxy]
```

Will start at the pods key, then since the next item is a [*] it assumes pods is a list, and will iterate over each item in the list to find where the key metadata.namespace is equal to grafana  It will then find the item where the key spec.containers is a list, and iterate over each item in the list to find where the key name is equal to istio-proxy 

Multiple filters can be added for a list, for example the above example could be modified to filter both by namespace and pod name:

```
pods.[metadata.namespace=grafana,metadata.name=operator].spec.containers.[name=istio-proxy]
```

To support map keys containing ".", [] syntax will also be used, e.g.,

```
namespaces.[metadata.namespace=grafana].metadata.labels.["some.key/label"]
```

Do not use the "[]" syntax for anything other than a key containing a ".", else the path will not be parsed correctly.

>[!IMPORTANT]
> The path will return only one item, the first item that matches the filters along the path. If no items match the filters, the path will return an empty map.

### Change Type Behavior

**Add**
* All keys in the path must exist, except for the last key. If you are trying to add a map, then use `value-map` and specify the existing root key.
* If a sequence is "added" to, then the value items will be appended to the sequence.

**Update**
* If a sequence is "updated", then the entire sequence will be replaced.

**Delete**
* Currently only supports deleting a key, error will be returned if the last item in the path resolves to a sequence.
* No values should be specified for delete.

## Executing Tests

Tests can be executed by specifying the `--run-tests` flag when running `lula dev validate`. E.g.,

```sh
lula dev validate -f ./validation.yaml --run-tests
```

This will execute the tests and print the test results to the console. 

To aid in debugging, the `--print-test-resources` flag can be used to print the resources used for each test to the validation directory, the filenames will be `<test-name>.json`.. E.g.,

```sh
lula dev validate -f ./validation.yaml --run-tests --print-test-resources
```

