# Compliance Evaluation

Evaluate serves as a method for verifying the compliance of a component/system against an established threshold to determine if it is more or less compliant than a previous assessment. 

## Expected Process

### No Existing Data

When no previous assessment exists, the initial assessment is made and stored with `lula validate`. This initial assessment by itself will always pass `lula evaluate` as there is no threshold for evaluation. Lula will automatically apply the `threshold` prop to the assessment result when writing the assessment result to a file that does not contain an existing assessment results artifact.

steps:
1. `lula validate`
2. `lula evaluate` -> Passes with no Threshold

### Existing Data (Intended Workflow)

In workflows run manually or with automation (such as CI/CD), there is an expectation that the threshold exists, and evaluate will perform an analysis of the compliance of the system/component against the established threshold.

steps:
1. `lula validate`
2. `lula evaluate` -> Passes or Fails based on threshold
