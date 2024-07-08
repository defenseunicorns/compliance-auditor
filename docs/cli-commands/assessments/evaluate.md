# Compliance Evaluation

Evaluate serves as a method for verifying the compliance of a component/system against an established threshold to determine if it is more or less compliant than a previous assessment. 

## Expected Process

### No Existing Data

When no previous assessment exists, the initial assessment is made and stored with `lula validate`. Lula will automatically apply the `threshold` prop to the assessment result when writing the assessment result to a file that does not contain an existing assessment results artifact. This initial assessment by itself will always pass `lula evaluate` as there is no threshold for evaluation, and the threshold prop with be set to `true`.

steps:
1. `lula validate -f component.yaml -o assessment-results.yaml`
2. `lula evaluate -f assessment-results.yaml` -> Passes with no Threshold -> Establishes Threshold

### Existing Data (Intended Workflow)

In workflows run manually or with automation (such as CI/CD), there is an expectation that the threshold exists, and evaluate will perform an analysis of the compliance of the system/component against the established threshold.

steps:
1. `lula validate -f component.yaml -o assessment-results.yaml`
2. `lula evaluate -f assessment-results.yaml` -> Passes or Fails based on threshold


## Scenarios for Consideration

Evaluate will determine which result is the threshold based on the following property:
```yaml
props:
  - name: threshold
    ns: https://docs.lula.dev/ns
    value: "true/false"
```

### Assessment Results Artifact

When evaluate is ran with a single assessment results artifact, it is expected that a single threshold with a `true` value exists. This will be identified and ran against the latest result to determine if compliance is less-than-equal (fail), equal (pass), or greater-than-equal (pass). When the comparison results in greater-than-equal, Lula will update the threshold `prop` for the latest result to `true` and set the previous result threshold prop to `false`.

### Comparing multiple assessment results artifacts

In the scenario where multiple assessment results artifacts are evaluated, there may be a multiple threshold results with a `true` value as Lula establishes a default `true` value when writing an assessment results artifact to a new file with no previous results present. In this case, Lula will use the older result as the threshold to determine compliance of the result.