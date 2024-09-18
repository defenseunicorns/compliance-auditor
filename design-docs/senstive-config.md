# Designing Sensitive Configuration

Author(s): @meganwolf0
Date Created: Sept 17, 2024
Status: DRAFT
Ticket: [#641](https://github.com/defenseunicorns/lula/issues/641)
Reviews Requested By: TBD

## Problem Statement

Given that validation data may incorporate sensitive information (e.g., API Keys, passwords, etc.), we need to determine the solution for 
1. Mapping those values to their source (e.g., environment variables, config files, secret store, etc.)
    -> I think as an initial cut, this should be scoped to variables only coming from the environment, maybe include functionality to source .env files
2. Ensuring that the values are masked in the output of the validation (e.g., when printing the component-definition)
3. Having the ability to manually template the values, specifically in the case of testing the validation content (probably in the purview of `lula dev`)

## Proposal

Comparing a few different options:

### Option 1: Go templates + string replacement
Basically, establish a prefix for the variables that are sensitive, then perform a string repace for the templated values -> also can apply a masking function when running the output?

(see https://github.com/defenseunicorns/lula/tree/517-go-template-testing for an example of templating something with "secret.xx")

### Option 2: HCL variable templating
I think this is probably generally more of a heavy-lift to get set-up, and will also introduce more verbosity when setting up variables, but could be a more flexible solution.

This is actually probably more relevant with generally templating, as it handles complex data structures...

## Scope and Requirements

IMO, the best way to define the scope is to determine which examples we are trying to cover:

1. Templating a sensitive API Key in a notional Domain API-Spec

```yaml
# ... rest of the validation ...
domain:
  type: api
  api-spec:
    requests:
    - name: local
      url: https://some.url/v1/api
      content-type: application/json
      headers:
        Authorization: Bearer {{ .secret.api_key }}
# ... rest of the validation ...
```

2. Templating a sensitive policy value...?

```yaml
# ...
# ...
```

## Implementation Details

Tiers of templating:

* Lula Tools Template -> Default templates constants + variables

Expand upon the implementation details here. Draw the reader’s attention to:

* Changes to existing systems.
* Creation of new systems.
* Impacts to the customer.
* Include code samples where possible.

## Metrics & Alerts

List any Metrics / Alerts that you plan to include in the system design

## Alternatives Considered

List any alternative solutions considered and how the proposed solution is a better fit.

## Non-Goals

List out anything that may be related to the solution, but won’t be covered by this solution.

## Future Improvements:

List out anything that won’t be included in this version of the feature/solution, but could be revisited or iterated upon in the future.

## Other Considerations:

List anything else that won’t be solved by the solution