# Resource Validation

Currently the Lula Kubernetes Domain and OPA provider operate under the current execution:
- The payload contains an array of ResourceRules
- The Kubernetes domain retrieves each of the applicable items from the cluster for each resource rule
- These items are collected into a single array of []unstructured.Unstructured
- The OPA provider converts this array into an []map[string]interface{}
  - Meaning `unstructured.Unstructured` -> `map[string]interface{}`
- The OPA provider then validates each item in the array against the rego policy


## Constraints

- This is an understood workflow - but it does not allow for the rich validation of many resources in a single validation.
  - IE In a single validation, I want to compare two resources to ensure the expected configuration matches
  - This could be done in two successive calls?
    - No, take for example the mapping of PeerAuthentications to Namespaces
      - We could have 1 -> N namespaces
      - For each namespace, we could have 0 -> N PeerAuthentications
      - if we wanted to validate that every applicable namespace has a PeerAuthentication AND every PeerAuthentication has mtls.mode set to STRICT
        - We need all namespaces
        - We need all PeerAuthentications
        - We need to ensure all namespaces have a PeerAuthentication
        - We need to ensure all PeerAuthentications have mtls.mode set to STRICT
- This is a very simple workflow - but it does not allow for the rich validation of many resources in a single validation.

## Multiple Resource Validation

- The payload contains an array of Resources
  - Each Resource contains a name and an array of ResourceRules (Adds an additional layer)
  - This allows us to collect multiple resources into groups for validation
- The Kubernetes domain retrieves each of the applicable items from the cluster for each resource rule
- These are returned as groups of items per Resource
- The OPA provider converts this array into an []map[string]interface{}
  - Meaning `[]unstructured.Unstructured` -> `map[string]interface{}`
  - Notice how this is now a collection of objects unlike previously
- The OPA provider then validates each item in the array against the rego policy

## How does this work?
- We need the ability to group resources into a collection 
  - Create a `resource` type
    - Each resource will contain some identifier and an array of resourcerules

-  