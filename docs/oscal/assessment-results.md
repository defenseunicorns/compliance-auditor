# Assessment Results

An [Assessment Result](https://pages.nist.gov/OSCAL/resources/concepts/layer/assessment/assessment-results/) is an OSCAL model to report on the specific assessment outcomes of a system. In Lula, the `validate` command creates an `assessment-result` object to enumerate the assessment of the input controls provided by the `component-definition`. These are reported as findings that are `satisfied` or `not-satisfied` as a result of the observations performed by the Lula validations.
```mermaid
flowchart TD
    A[Assessment Results]-->|compose|C[Finding 1]
    A[Assessment Results]-->|compose|G[Finding 2]
    B(Control)-->|satisfied by|C
    B(Control)-->|satisfied by|G
    C -->|compose|D[Observation 1]
    C -->|compose|E[Observation 2]
    C -->|compose|F[Observation 3]
    G -->|compose|F[Observation 3]
    G -->|compose|H[Observation 4]
    D --> I(Lula Validation A)
    E --> J(Lula Validation B)
    F --> K(Lula Validation C)
    H --> L(Lula Validation D)
```

## Observation Results
Based on the structure outlined, the results of the observations impact the findings, which in turn result in the decision for the control as `satisfied` or `not-satisfied`. The observations are aggregated to the findings as `and` operations, such that if a single observation is `not-satisfied` then the associated finding is marked as `not-satisfied`.

The way Lula performs evaluations default to a conservative reporting of a `not-satisfied` observation. The only `satisfied` observations occur when a domain provides resources and those resources are evaluated by the policy such that the policy will pass. If a Lula Validation [cannot be evaluated](#not-satisfied-conditions) then it will by default return a `not-satisfied` result.

### Not-satisfied conditions
The following conditions enumerate when the Lula Validation will result in a `not-satisfied` evaluation. These cases exclude the case where the Lula validation policy has been evaluated and returned a failure.
- Malformed Lula validation -> bad validation structure
- Missing resources -> No resources are found as input to the policy
- Missing reference -> If a remote or local reference is invalid
- Executable validations disallowed -> If a validation is executable but has not been allowed to run

## Structure
The primary structure for Lula production and operation of `assessment-results` for determinism is as follows:
- Results are sorted by `start` time in descending order
- Findings are sorted by `target.target-id` in ascending order
- Observations are sorted by `collected` time in ascending order
- Back Matter Resources are sorted by `title` in ascending order.
