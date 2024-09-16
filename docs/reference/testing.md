# Testing

Testing is a key part of Lula Validation development. Since the results of the Lula Validations are determined by the policy set by the `provider`, those policies must be tested to ensure they are working as expected.

## Validation Testing

In the Lula Validation, a `tests` property is used to specify the each test that should be performed against the validation. Each test is a map of the following properties:

- `name`: The name of the test
- `permutations`: An array of permutations to be performed on the validation
- `expected-result`: The expected result of the test - pass or fail

### Permutations

A permutation is a map of the following properties:

- `path`: The path to the resource to be modified. The path syntax is described below.
- `type`: The type of operation to be performed on the resource
    - `update`: (default) updates the resource with the specified value
    - `delete`: deletes the field specified
    - `add`: adds the specified value
- `value`: The value to be used for the operation (map[string]interface{} or string)

### Path Syntax

This feature uses the kyaml library to inject data into the resources, so the path syntax is based on this library. 

The path should be a "." delimited string that specifies the keys along the path to the resource seeking to be modfied. In addition to keys, a list can be specified by using the `[]` syntax. For example, the following path:

```yaml
pods.[metadata.namespace=grafana].spec.containers.[name=istio-proxy]
```

Will start at the `pods` key, then since the next item is a [] it assumes `pods` is a list, and will iterate over each item in the list to find where the key `metadata.namespace` is equal to `grafana`. It will then find the item where the key `spec.containers` is a list, and iterate over each item in the list to find where the key `name` is equal to `istio-proxy`.

The path will return only one item, the first item that matches the filters along the path. If no items match the filters, the path will return an empty map.

## Examples
See `src/test/unit/types/validation-all-pods.yaml` for an exmple of a validation with tests.
