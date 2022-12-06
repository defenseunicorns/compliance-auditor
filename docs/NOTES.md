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