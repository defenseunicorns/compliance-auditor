# Validation

The Assessment Results OSCAL model stores the `observation` list that relate to a specific `finding` (`finding.related_observations`). These findings are the control evaluations of the component/system, more detail can be found in [Assessment Results](../assessment-results.md). The observations are generated by Lula during the execution of a `validate` operation, which is underpinned by multiple executions of the [Lula Validations](../../reference/README.md).

To map these obserservations to the specific Lula Validation, the `validation` prop is used to identify the Lula Validation by unique identifier.

## Example

After the `validate` operation where `observation` are generated in the `assessment-results` - Lula will add a `observation.props` entry to the result in the following format:
```yaml
props:
  - name: validation
    ns: https://docs.lula.dev/oscal/ns
    value: "#8894c5cd-e27b-437e-a146-b6b4e9f2be78"
```

The `value` here is the Lula Validation's UUID.
