# Testing

Testing is a key part of Lula Validation development. Since the results of the Lula Validations are determined by the policy set by the `provider`, those policies must be tested to ensure they are working as expected.

## Validation Testing

In the Lula Validation, a `tests` property is used to specify each test that should be performed against the validation. Each test is a map of the following properties:

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

Will start at the pods key, then since the next item is a [*=*] it assumes pods is a list, and will iterate over each item in the list to find where the key `metadata.namespace` is equal to `grafana`. It will then find the `containers` list item in `spec`, and iterate over each item in the list to find where the key `name` is equal to `istio-proxy`.

Multiple filters can be added for a list, for example the above example could be modified to filter both by namespace and pod name:

```
pods.[metadata.namespace=grafana,metadata.name=operator].spec.containers.[name=istio-proxy]
```

To support map keys containing ".", [] syntax will also be used, e.g.,

```
namespaces.[metadata.namespace=grafana].metadata.labels.["some.key/label"]
```

Additionally, individual list items can be found via their index, e.g.,

```
namespaces.[0].metadata.labels
```

Which will point to the labels key of the first namespace. Additionally, a `[-]` can be used to specify the last item in the list.

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

A note about replacing a key with an empty map - due to the way the `kyaml` library works, simply trying to overwrite an existing key with an empty map will not yield a removal of all the existing data of the map, it will just try and merge the differences, which is possibly not the desired outcome. To replace a map with an empty map, you must combine `delete` a change type and `add` a change type, e.g.,

```yaml
changes:
  - path: pods.[metadata.namespace=grafana].metadata.labels
    type: delete
  - path: pods.[metadata.namespace=grafana].metadata
    type: add
    value-map: 
      labels: {}
```

Which will delete the existing labels map and then add an empty map, such that the "labels" key will still exist but will be an empty map.

## Executing Tests

Tests can be executed by specifying the `--run-tests` flag when running both `lula validate` and `lula dev validate`, however the output of either will be slightly different.

### lula validate
When running `lula validate ... --run-tests`, a test results summary will be printed to the console, while the test results yaml file will be written to the same directory as the output data (either default to directory of the `component-definition` source file or to the directory specified by the `--output-file` flag).

E.g., Running validate on a component-definition with two validations, one with tests and one without:
```sh
lula validate -f ./component.yaml --run-tests
```

Will print the test results summary to the console as:
```sh
Test Results: 1 passing, 0 failing, 1 missing
```

And will print a test results yaml file to the same directory as the output data:
```yaml
61ec8808-f0f4-4b35-9a5b-4d7516053534:
    name: test-validation
    test-results: []
82099492-0601-4287-a2d1-cc94c49dca9b:
    name: test-validation-with-tests
    test-results:
        - test-name: change-image-name
          pass: true
          result: not-satisfied
        - test-name: no-containers
          pass: true
          result: not-satisfied
```
> Note that `61ec8808-f0f4-4b35-9a5b-4d7516053534` is the UUID of the validation without tests, and `82099492-0601-4287-a2d1-cc94c49dca9b` is the UUID of the validation with tests.

### lula dev validate
When executing `lula dev validate ... --run-tests`, the test results data will be written directly to console.

E.g., Running dev validate on a Lula validation with two tests:
```sh
lula dev validate -f ./validation.yaml --run-tests
```

Will print the test results to the console as:
```sh
  ✔  Pass: change-image-name
  •  Result: not-satisfied
  ✔  Pass: no-containers
  •  Result: not-satisfied
```

To aid in debugging, the `--print-test-resources` flag can be used to print the resources used for each test to the validation directory, the filenames will be `<test-name>.json`.. E.g.,

```sh
lula dev validate -f ./validation.yaml --run-tests --print-test-resources
```
