# Develop a Validation

This document describes some best practices for developing and using a [Lula Validation](), the primary mechanism for evaluation of system's compliance against a specified control.

## About

Lula Validations are constructed by establishing the collection of measurements about a system, given by the specified **Domain**, and the evaluation of adherence, performed by the **Provider**. 

The currently supported domains are:
- API
- Kubernetes

The currently supported providers are:
- OPA (Open Policy Agent)
- Kyverno

## Creating a Sample Validation

Here, we will step through creating a sample validation using the Kubernetes domain and OPA provider. Generating a validation is in the scope of answering some control or standard. For instance, our control might be something like "system implements test application as target for development purposes". Our validation should then seek to _prove_ that some "test application" is running in our domain, i.e., Kubernetes.

### Pre-Requistes

* Lula installed
* Kubectl
* Helm
* A running Kubernetes cluster
    - Kind
        - `kind create cluster -n lula-test`
    - K3d
        - `k3d cluster create lula-test`
* Podinfo deployed in the cluster
    ```sh
    $ helm repo add podinfo https://stefanprodan.github.io/podinfo

    $ helm upgrade -i my-release podinfo/podinfo -n podinfo --create-namespace
    ```

### Steps

>[!NOTE]
> Demo files can be found in the lula repository under `demo/develop-validation`

1. Assume we have some component definition for Podinfo with the associated standard we are trying to prove the system satisfies:
    ```yaml
    component-definition:
      uuid: a506014d-cb8a-4db9-ac48-ef72f7209a60
      metadata:
        last-modified: 2024-07-11T13:38:09.633174-04:00
        oscal-version: 1.1.2
        published: 2024-07-11T13:38:09.633174-04:00
        remarks: Lula Generated Component Definition
        title: Component Title
        version: 0.0.1
      components:
      - uuid: 75859c1e-30f5-4fde-9ad4-c79f863b049f
        type: software
        title: podinfo
        description: Sample application
        control-implementations:
        - uuid: a3039927-839c-5745-ac4e-a9993bcd60ed
          source: https://github.com/defenseunicorns/lula
          description: Control Implementation Description
          implemented-requirements:
          - uuid: 257d2b2a-fda7-49c5-9a2b-acdc995bc8e5
            control-id: ID-1
            description: >-
              Podinfo, a sample application, is deployed into the cluster and exposed for testing purposes.
            remarks: >-
              System implements test application as target for development purposes.
    ```

2. We recognize that we can satisfy this control by proving that podinfo is alive in the cluster. If we know nothing about podinfo, we may first want to identify which Kubernetes constructs are used in it's configuration:
    ```sh
    $ kubectl get all -n podinfo 

    NAME                                     READY   STATUS    RESTARTS   AGE
    pod/my-release-podinfo-fb6d4888f-ptlss   1/1     Running   0          17m

    NAME                         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)             AGE
    service/my-release-podinfo   ClusterIP   10.43.172.65   <none>        9898/TCP,9999/TCP   17m

    NAME                                 READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/my-release-podinfo   1/1     1            1           17m

    NAME                                           DESIRED   CURRENT   READY   AGE
    replicaset.apps/my-release-podinfo-fb6d4888f   1         1         1       17m
    ```

3. Now that we know what resources are in the `podinfo` namespace, we can use our kubernetes knowledge to deduce that proving podinfo is healthy in the cluster could be performed by looking at the `status` of the podinfo deployment for the `replicas` value to match `readyReplicas`:
    ```sh
    $ kubectl get deployment my-release-podinfo -n podinfo -o json | jq '.status'

    {
    "availableReplicas": 1,
    "conditions": [
      {
        "lastTransitionTime": "2024-07-11T17:36:53Z",
        "lastUpdateTime": "2024-07-11T17:36:53Z",
        "message": "Deployment has minimum availability.",
        "reason": "MinimumReplicasAvailable",
        "status": "True",
        "type": "Available"
      },
      {
        "lastTransitionTime": "2024-07-11T17:36:53Z",
        "lastUpdateTime": "2024-07-11T17:36:56Z",
        "message": "ReplicaSet \"my-release-podinfo-fb6d4888f\" has successfully progressed.",
        "reason": "NewReplicaSetAvailable",
        "status": "True",
        "type": "Progressing"
      }
    ],
    "observedGeneration": 1,
    "readyReplicas": 1,
    "replicas": 1,
    "updatedReplicas": 1
    }
    ```

4. With this we should now have enough information to write our Lula Validation! First construct the top-matter `metadata`:

    Run `lula tools uuidgen` to get a unique ID for your validation
    ```bash
    $ lula tools uuidgen              
    ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
    ```

    Add a `validation.yaml` file with the following
    ```yaml
    metadata:
      name: check-podinfo-health
      uuid: ad38ef57-99f6-4ac6-862e-e0bc9f55eebe
    ```

5. Construct the `domain`:

    Since we are extracting Kubernetes manifest data as validation "proof", the domain we use should be `kubernetes`.
    ```yaml
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
    ```
    
    Note a few things about the specification for obtaining these kubernetes resources:
    - `resources` key is used as an array of resources we are asking for from the cluster
        - `name` is the keyword that will be used as an input to the policy, stated below in the provider. Note - to play nicely with the policy, it is best to make this a single word, camel-cased if desired.
        - `resource-rule` is the api specification for the resource being extracted
            - `name` is the name of our deployment, `my-release-podinfo`
            - `namespaces` is the list of namespaces, one can be provided but must be in list format
            - `group`, `version`, `resource` is the compliant values to access the [kubernetes API](https://pkg.go.dev/k8s.io/apimachinery@v0.30.2/pkg/runtime/schema)

    See [reference](../reference/README.md) for more information about the Lula Validation schema and kubernetes domain.

6. Construct the `provider` and write the OPA policy:

    Any provider should be compatible with the domain outputs, here we've decided to use OPA and rego, so our `provider` section is as follows:
    ```yaml
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
            input.podinfoDeployment.status.replicas > 0
            input.podinfoDeployment.status.availableReplicas == input.podinfoDeployment.status.replicas
            msg := "Number of replicas > 0 and all replicas are available."
          } else = {"result": false, "msg": msg} {
            msg := "Podinfo not available."
          }
        output:
          validation: validate.validate
          observations:
            - validate.msg
    ```

    The Rego policy language can be a little funny looking at first glance, check out both the rego docs and the [OPA Provider](../reference/providers/opa-provider.md) reference for more information about rego.

    With that said, some things are important to highlight about the policy
    - `package validate` is mandatory at the top (you can use any package name you want, but if a different value is used the `output.validation` needs to be updated accordingly)
    - `import rego.v1` is optional, but recommended as OPA looks to [upgrade to v1](https://www.openpolicyagent.org/docs/latest/opa-1/)
    - The "Default values" section is best practice to set these to protect against a result that yields undefined values for these variables
    - The "Validation result" section defines the rego evaluation on the `input.podinfoDeployment` - checking that both the number of replicas is greater than 0 and the available and requested replicas are equal.

7. Putting it all together we are left with the `validation.yaml`, let's run some commands to validate our validation:

    Get the resources to visually inspect that they are what you expect from the `domain` and in the right struction for the provider's policy:
    ```sh
    $ lula dev get-resources -f validation.yaml -o resources.json
    ```

    The result should be a `resources.json` file that looks roughly as follows:
    ```json
    {
      "podinfoDeployment": {
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        # ... rest of the json
      }
    }
    ```

    Now check the validation is resulting in the expected outcome:
    ```sh
    $ lula dev validate -f validation.yaml                        
    •  Observations:
    •  --> validate.msg: Number of replicas > 0 and all replicas are available.
    •  Validation completed with 1 passing and 0 failing results
    ```
    
    If we expected this validation to fail, we would have added `-e=false`

8. Now that we have our baseline validation, and we know it is returning an expected result for our current cluster configuration, we should probably ensure that the policy results are successful when other resource cases exist. There are really two options here:
    
    1) Manually modify the resources in your cluster and re-run `lula dev validate`

    2) Manually modify the `resources.json` and test those

        If we have a test cluster, perhaps changing some things about it is acceptable, but for this case I'm just going to take the path of least resistance and modify the `resources.json`:

        Copy your `resources.json` and rename to `resources-bad.json`. First, find `podinfoDeployment.status.replicas` and change the value to 0. Run `lula dev validate` with those resources as the input, along with our expected failure outcome:

        ```sh
        $ lula dev validate -f validation.yaml -r resources-bad.json -e=false                    
        • Observations: 
        •  --> validate.msg: Podinfo not available.                                                              
        •  Validation completed with 0 passing and 1 failing results 
        ```

        Success! Additional conditions can be tested this way to fully stress-test the validity of the policy.

9. Finally, we can bring this back into the `component-definition`. This validation should be added as a link to the respective  `implemented-requirement`:
    ```yaml
    # ... rest of component definition
        implemented-requirements:
          - uuid: 257d2b2a-fda7-49c5-9a2b-acdc995bc8e5
            control-id: ID-1
            description: >-
            Podinfo, a sample application, is deployed into the cluster and exposed for testing purposes.
            remarks: >-
            System implements test application as target for development purposes.
            links:
              - href: 'file:./validation.yaml
                ref: lula
                text: Check that Podinfo is healthy
    ```

10. Now that we have our full OSCAL Component Definition model specified, we can take this off to `validate` and `evaluate` the system!

## Limitations

We are aware that many of these validations are brittle to environment changes, for instance if namespaces change. This is a known limitation and on our roadmap as something to offer a possible templating solution for in the future.

Additionally, since we are adding these validations to OSCAL yaml documents, there is some ugliness with having to compose strings of yaml into yaml. We support "remote" validations, where instead of a reference to a backmatter uuid, instead a link to a file is provided. A limitation of that currently is that it does not support authentication if the remote link is in a protected location. 