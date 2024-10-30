# Testing

Testing is a key part of Lula Validation development. Since the results of the Lula Validations are determined by the policy set by the `provider`, those policies must be tested to ensure they are working as expected.

## Validation Testing

In the Lula Validation, a `tests` property is used to specify the each test that should be performed against the validation. Each test is a map of the following properties:

- `name`: The name of the test
- `changes`: An array of changes or transformations to be applied to the resources used in the test validation
- `expected-result`: The expected result of the test - pass or fail

### Change

A change is a map of the following properties:

- `path`: The path to the resource to be modified. The path syntax is described below.
- `type`: The type of operation to be performed on the resource
    - `update`: (default) updates the resource with the specified value
    - `delete`: deletes the field specified
    - `add`: adds the specified value
- `value`: The value to be used for the operation (string)
- `value-map`: The value to be used for the operation (map[string]interface{})

#### Path Syntax

This feature uses the kyaml library to inject data into the resources, so the path syntax is based on this library. 

The path should be a "." delimited string that specifies the keys along the path to the resource seeking to be modified. In addition to keys, a list can be specified by using the “[]” syntax. For example, the following path:

```
pods.[metadata.namespace=grafana].spec.containers.[name=istio-proxy]
```

Will start at the pods key, then since the next item is a [*] it assumes pods is a list, and will iterate over each item in the list to find where the key metadata.namespace is equal to grafana  It will then find the item where the key spec.containers is a list, and iterate over each item in the list to find where the key name is equal to istio-proxy 

Multiple filters can be added for a list, for example the above example could be modified to filter both by namespace and pod name:

```
pods.[metadata.namespace=grafana, metadata.name=operator].spec.containers.[name=istio-proxy]
```

To support map look-ups, [] will also be used, but when NOT separated by a “.” the item from the map will just be identified directly, e.g.,

```
namespaces.[metadata.namespace=grafana].metadata.labels["istio-injection"]
```

Note: The path will return only one item, the first item that matches the filters along the path. If no items match the filters, the path will return an empty map.

#### Change Type Behavior

**Add**
* All keys in the path must exist, except for the last key. If you are trying to add a map, then use `value-map` and specify the existing root key.
* If a sequence is "added" to, then the value items will be appended to the sequence.

**Update**
* If a sequence is "updated", then the entire sequence will be replaced.

**Delete**
* Currently only supports deleting a key, error will be returned if the last item in the path resolves to a sequence.
* No values should be specified for delete.

## Examples
See `src/test/unit/types/validation-all-pods.yaml` for an exmple of a validation with tests.
