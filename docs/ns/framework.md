# Framework

Component Definitions are built to be modular / re-usable components that can be applied against one-to-many systems. As such, it is expected that a single Component Definition can contain one-to-many Components with each Component containing one-to-many Control Implementations.

Each Control Implementation may reference a source to many disparate and connected or non-connected standards. As part of the Lula `validate` operation, the `framework` OSCAL prop exists for allowing a defined validation activity to include multiple Control Implementations into a single result for continued analysis. 

## Example

This use of this prop is defined in `control-implementation[_].props` in the following format:

```yaml
props:
  - name: framework
    ns: https://docs.lula.dev/ns
    value: impact-level-x
```

This example would be for the specification of the "Impact Level X" collection of Control Implementations.
