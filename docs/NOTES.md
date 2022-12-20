# Notes

## Extension
- Existence of the tool/resource(s)
    - Can we fix this issue in the Kyverno CLI?
    - Or do we provide a different layer for providing this validation?

## Kyverno

### Limitations
- Wildcard for match any "kind" does not work as specified
- Cannot correlate between multiple resources well
    - Example: Given global context - cannot verify that all namespaces include a particular resource (CRD existence may be important)
    - Review the use of [External Data Sources](https://kyverno.io/docs/writing-policies/external-data-sources/) through the CLI
        - There are known issues here - but if resolvable could provide a mechanism for reporting on the existence of a resource
        - This would also allow establishing 1 -> Many relationships among resources