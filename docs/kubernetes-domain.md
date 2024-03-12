# Kubernetes Domain

The Kubernetes domain provides Lula with a common interface for data collection of Kubernetes artifacts for use across many Lula Providers. 

## Payload Expectation

The validation performed when using the Kubernetes domain is as follows:

```yaml
resources:
- name: podsvt                      # Required - Identifier for use in the rego below
  resourceRule:                     # Required - resource selection criteria, at least one resource rule is required
    Name:                           # Optional - Used to retrieve a specific resource in a single namespace
    Group:                          # Required - empty or "" for core group
    Version: v1                     # Required - Version of resource
    Resource: pods                  # Required - Resource type
    Namespaces: [validation-test]   # Required - Namespaces to validate the above resources in. Empty or "" for all namespace pr non-namespaced resources
    Field:                          # Optional - Field to grab in a resource if it is in an unusable type, e.g., string json data. Must specify named resource to use.
      jsonpath:                     # Required - Jsonpath specifier of where to find the field from the top level object
      type:                         # Optional - Accepts "json" or "yaml". Default is "json".
      base64:                       # Optional - Boolean whether field is base64 encoded

```

> [!Tip]
> Lula supports eventual-consistency through use of an optional `wait` field. 

```yaml
wait:
  condition: Ready
  kind: pod/test-pod-wait
  namespace: validation-test
  timeout: 30s
resources:
- name: podsvt
  resourceRule:
    Group:
    Version: v1
    Resource: pods
    Namespaces: [validation-test]
```

## Lists vs Named Resource

When Lula retrieves all targeted resources (bounded by namespace when applicable), the payload is a list of resources. When a resource Name is specified - with a target Namespace - the payload will be a single object. 

### Example

Let's get all pods in the `validation-test` namespace and evaluate them with the OPA provider:
```yaml
target:
  provider: opa
  domain: kubernetes
  payload:
    resources:
    - name: podsvt
      resourceRule:
        Group:
        Version: v1
        Resource: pods
        Namespaces: [validation-test]
    rego: |
      package validate

      import future.keywords.every

      validate {
        every pod in input.podsvt {
          podLabel := pod.metadata.labels.foo
          podLabel == "bar"
        }
      }
```

> [!IMPORTANT]
> Note how the payload contains a list of items that can be iterated over. The `podsvt` field is the name of the field in the payload that contains the list of items.

Now let's retrieve a single pod from the `validation-test` namespace:

```yaml
target:
  provider: opa
  domain: kubernetes
  payload:
    resources:
    - name: podvt
      resourceRule:
        Name: test-pod-label
        Group:
        Version: v1
        Resource: pods
        Namespaces: [validation-test]
    rego: |
      package validate

      validate {
        podLabel := input.podvt.metadata.labels.foo
        podLabel == "bar"
      }
```

> [!IMPORTANT]
> Note how the payload now contains a single object called `podvt`. This is the name of the resource that is being validated.

## Extracting Resource Field Data
Many of the tool-specific configuration data is stored as json or yaml text inside configmaps and secrets. Some valuable data may also be stored in json or yaml strings in other resource locations, such as annotations. The "Field" parameter of the "ResourceRule" allows this data to be extracted and used by the Rego.

Here's an example of extracting `config.yaml` from a test configmap:
```yaml
target:
  provider: opa
  domain: kubernetes
  payload:
    resources:
    - name: configdata
      resourceRule:
        Name: test-configmap
        Group:
        Version: v1
        Resource: pods
        Namespaces: [validation-test]
        Field:
          jsonpath: .data.my-config.yaml
          type: yaml
    rego: |
      package validate

      validate {
        configdata.configuration.foo == "bar"
      }
```

Where the raw ConfigMap data would look as follows:
```yaml
configuration:
  foo: bar
  anotherValue:
    subValue: ba
```