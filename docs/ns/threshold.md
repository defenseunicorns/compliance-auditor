# Threshold

The Assessment Results OSCAL model supports the storage of many Result objects pertaining to the assessment of a component/system. 

Each of these Results may establish a level of Compliance that indicates how compliant a component/system was at any point in time (typically during assessment). Lula leverages the Assessment Results model to store results of each `validate` operation while maintaining a `threshold` indicating when a component/system was most compliant. 

This field is automatically maintained as Lula processes the identification of a `threshold` and updating it as required when a component/system becomes more compliant that a previous threshold. 

## Example

After the initial `validate` operation - Lula will add a `results[_].props` entry to the result in the following format:
```yaml
props:
  - name: threshold
    ns: https://docs.lula.dev/ns
    value: "false"
```

As indicated in the [evaluate section](../cli-commands/assessments/evaluate.md), it is expected that `lula evaluate` be executed with an Assessment Results artifact containing more than one comparable result. In the event a single result exists, Lula will automatically add a `threshold` property to the result and set the `value` to `true`.