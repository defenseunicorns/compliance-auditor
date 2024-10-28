# Kubernetes Domain

The Kubernetes domain provides Lula access to data collection of Kubernetes resources.

>[!Important]
>This domain supports both read and write operations on the Kubernetes cluster in the current context, so use with care

## Specification

The Kubernetes specification can be used to return Kubernetes resources as JSON data to the providers. The specification allows for both manifest and resource field data to be retreived.

The following sample `kubernetes-spec` yields a list of all pods in the `validation-test` namespace.

```yaml
domain:
  type: kubernetes
  kubernetes-spec:
    resources:                          # Optional - Group of resources to read from Kubernetes
    - name: podsvt                      # Required - Identifier to the list or set read by the policy
      resource-rule:                    # Required - Resource selection criteria, at least one resource rule is required
        name:                           # Optional - Used to retrieve a specific resource in a single namespace
        group:                          # Optional - empty or "" for core group
        version: v1                     # Required - Version of resource
        resource: pods                  # Required - Resource type (API-recognized type, not Kind)
        namespaces: [validation-test]   # Optional - Namespaces to validate the above resources in. Empty or "" for all namespace or non-namespaced resources
        field:                          # Optional - Field to grab in a resource if it is in an unusable type, e.g., string json data. Must specify named resource to use.
          jsonpath:                     # Required - Jsonpath specifier of where to find the field from the top level object
          type:                         # Optional - Accepts "json" or "yaml". Default is "json".
          base64:                       # Optional - Boolean whether field is base64 encoded
```

Lula supports eventual-consistency through use of an optional `wait` field in the `kubernetes-spec`. This parameter supports waiting for a specified resource to be `Ready` in the cluster. This may be particularly useful if evaluating the status of a selected resource or evaluating the children of a specified resource.

```yaml
domain:
  type: kubernetes
  kubernetes-spec:
    wait:                               # Optional - Group of resources to read from Kubernetes
      group:                            # Optional - Empty or "" for core group
      version: v1                       # Required - Version of resource
      resource: pods                    # Required - Resource type (API-recognized type, not Kind)Required - Resource type (API-recognized type, not Kind)
      name: test-pod-wait               # Required - Name of the resource to wait for
      namespace: validation-test        # Optional - For namespaced resources
      timeout: 30s                      # Optional - Defaults to 30s
    resources:
    - name: podsvt
      resource-rule:
        group:
        version: v1
        resource: pods
        namespaces: [validation-test]
```

> [!Tip]
> Both `resources` and `wait` use the Group, Version, Resource constructs to identify the resource to be evaluated. To identify those using `kubectl`, executing `kubectl explain <resource/kind/short name>` will provide the Group and Version, the `resource` field is the API-recognized type and can be confirmed by consulting the list provided by `kubectl api-resources`.

### Resource Creation

The Kubernetes domain also supports creating, reading, and destroying test resources in the cluster. This feature should be used with caution since it's writing to the cluster and ideally should be implemented on separate namespaces to make any erroneous outcomes easier to mitigate.

```yaml
domain:
  type: kubernetes
  kubernetes-spec:
    create-resources:                   # Optional - Group of resources to be created/read/destroyed in Kubernetes
      - name: testPod                   # Required - Identifier to be read by the policy
        namespace: validation-test      # Optional - Namespace to be created if applicable (no need to specify if ns exists OR resource is non-namespaced)
        manifest: |                     # Optional - Manifest string for resource(s) to create; Only optional if file is not specified
          <some manifest(s)>
        file: '<some url>'              # Optional - File name where resource(s) to create are stored; Only optional if manifest is not specified. Currently does not support relative paths.
```

In addition to simply creating and reading individual resources, you can create a resource, wait for it to be ready, then read the possible children resources that should be created. For example the following `kubernetes-spec` will create a deployment, wait for it to be ready, and then read the pods that should be children of that deployment:

```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    create-resources:
    - name: testDeploy
      manifest: |
        apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: test-deployment
          namespace: validation-test
        spec:
          replicas: 1
          selector:
            matchLabels:
              app: test-app
          template:
            metadata:
              labels:
                app: test-app
            spec:
              containers:
              - name: test-container
                image: nginx
    wait:
      group: apps
      version: v1
      resource: deployments
      name: test-deployment
      namespace: validation-test
    resources:
    - name: validationTestPods
      resource-rule:
        version: v1
        resource: pods
        namespaces: [validation-test]
```

> [!NOTE]
> The `create-resources` is evaluated prior to the `wait`, and `wait` is evaluated prior to the `resources`.

## Lists vs Named Resource

When Lula retrieves all targeted resources (bounded by namespace when applicable), the payload is a list of resources. When a resource Name is specified - the payload will be a single object. 

### Example

Let's get all pods in the `validation-test` namespace and evaluate them with the OPA provider:
```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: podsvt
      resource-rule:
        group:
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
```

> [!IMPORTANT]
> Note how the rego evaluates a list of items that can be iterated over. The `podsvt` field is the name of the field in the kubernetes-spec.resources that contains the list of items.

Now let's retrieve a single pod from the `validation-test` namespace:

```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: podvt
      resource-rule:
        name: test-pod-label
        group:
        version: v1
        resource: pods
        namespaces: [validation-test]
provider: 
  type: opa
  opa-spec:  
    rego: |
      package validate

      validate {
        podLabel := input.podvt.metadata.labels.foo
        podLabel == "bar"
      }
```

> [!IMPORTANT]
> Note how the rego now evaluates a single object called `podvt`. This is the name of the resource that is being validated.

We can also retrieve a single cluster-scoped resource as follows, where the rego evaluates a single object called `namespaceVt`.

```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: namespaceVt
      resource-rule:
        name: validation-test
        version: v1
        resource: namespaces
provider: 
  type: opa
  opa-spec:  
    rego: |
      package validate

      validate {
        input.namespaceVt.metadata.name == "validation-test"
      }
```

## Extracting Resource Field Data
Many of the tool-specific configuration data is stored as json or yaml text inside configmaps and secrets. Some valuable data may also be stored in json or yaml strings in other resource locations, such as annotations. The `field` parameter of the `resource-rule` allows this data to be extracted and used by the Rego.

Here's an example of extracting `config.yaml` from a test configmap:
```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: configdata
      resource-rule:
        name: test-configmap
        group:
        version: v1
        resource: configmaps
        namespaces: [validation-test]
        field:
          jsonpath: .data.my-config.yaml
          type: yaml
provider: 
  type: opa
  opa-spec:
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
This field also supports grabbing secret data using the optional `base64` field as follows
```yaml
domain: 
  type: kubernetes
  kubernetes-spec:
    resources:
    - name: secretdata
      resource-rule:
        name: test-secret
        group:
        version: v1
        resource: secrets
        namespaces: [validation-test]
        field: 
          jsonpath: .data.secret-value.json
          type: json
          base64: true
```
