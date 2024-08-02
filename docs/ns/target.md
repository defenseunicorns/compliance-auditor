# Target

The `target` prop is a prop used by Lula to automate identification of a result back to a given `Control Implementation` source or a collection of multiple sources when the [framework prop](./framework.md) is set. 

This is used by Lula to allow for `lula validate` and `lula evaluate` operations to target a specific standard or collection of standards.

## Example

This prop can be identified when reviewing `Assessment Results` artifacts for each `result[_].props` as follows:

```yaml
props:
  - name: target
    ns: https://docs.lula.dev/ns
    value:  https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.yaml
```
