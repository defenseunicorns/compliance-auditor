# System Security Plan

**NOTE:** This Document is in an active research phase.

## System Security Plan Generate

To generate a system security plan, you need the following context:
- The component definition `--component`
- The profile  source `--profile` or catalog source `--catalog`; IE `https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.json`

The following command will generate a system security plan with the above context:

```bash

lula generate system-security-plan --component .src/test/unit/valid-component.yaml --catalog https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.json

```

There are optional flags that can be added to the command to generate a system security plan:

- The output file of the component `-o` or `--output`; `oscal-system-security-plan.yaml`

## System Security Plan Generate Context

The `system-security-plan` can be generated using the upstream catalog and/or profile in conjunction with the `component-definition`. There are net new fields that are apart of the `system-security-plan` that are not within the `component-definition` or catalog/profile that currently do not make sense to add as props. Those items are under the section `Elements in SSP Not in Component Definition`. There are items that are not in the `system-security-plan` but also not in the `component-definition` that currently do make sense to create as props. Those items are under the section `Elements NOT in Component Definition that need added for SSP Generate`. Lastly as a note there are items within the `component-definition` that are not used in the `system-security-plan` that can be found under the section Elements NOT in `Component Definition that need added for SSP Generate`.

The items in `Elements in SSP Not in Component Definition` need further context to fill in the missing elements as well as establish data across OSCAL models. Some examples of the data fields are within the `metadata` fields such as `responsible-roles`, `responsible-parties`, and `parties` that can be added to the `system-security-plan` that do not directly map from the `component-definition` field. Additional context can be added through common OSCAL fields such as `props`, `links`, and `remarks`.

### Further Research Fields

The following fields need further research to further enhance generating an SSP.

- `inherited`
- `export`

### Elements in Component Definition Not in SSP

- `import-component-definitions`
- `capabilities`
  - `uuid`
  - `name`
  - `description`
  - `props`
  - `links`
  - `incorporates-components`
    - `component-uuid`
    - `description`

### Elements in SSP Not in Component Definition

- `system-characteristics`
  - `system-ids`
    - `identifier-type`
    - `id`
  - `system-name`
  - `system-name-short`
  - `description`
  - `security-sensitivity-level`
  - `system-information`
    - `information-types`
      - `id`
      - `title`
      - `description`
      - `security-objective-confidentiality`
      - `security-objective-integrity`
      - `security-objective-availability`
  - `security-impact-level`
    - `security-objective-confidentiality`
    - `security-objective-integrity`
    - `security-objective-availability`
  - `status`
    - `state`
    - `remarks`
  - `authorized-boundary`
    - `description`
    - `props`
    - `links`
    - `diagrams`
      - `uuid`
      - `description`
      - `props`
      - `links`
      - `caption`
      - `remarks`
    - `remarks`
  - `network-architecture`
    - `description`
    - `props`
    - `links`
    - `diagrams`
      - `uuid`
      - `description`
      - `props`
      - `links`
      - `caption`
      - `remarks`
    - `remarks`
  - `data-flow`
    - `description`
    - `props`
    - `links`
    - `diagrams`
      - `uuid`
      - `description`
      - `props`
      - `links`
      - `caption`
      - `remarks`
    - `remarks`
  - `props`
  - `links`
  - `remarks`
- `system-implementation`
  - `users`
    - `uuid`
    - `title`
    - `short-name`
    - `description`
    - `props`
    - `links`
    - `role-ids`
    - `authorized-privileges`
      - `title`
      - `description`
      - `functions-performed`
    - `remarks`
  - `leveraged-authorizations`
    - `uuid`
    - `title`
    - `props`
    - `links`
    - `party-uuid`
    - `date-authorized`
    - `remarks`
  - `inventory-items`
    - `uuid`
    - `description`
    - `props`
    - `links`
    - `remarks`
- `control-implementation`
  - `implemented-requirements`
    - `statements`
      - `satisfied`
        - `uuid`
        - `responsibility`
        - `description`
        - `props`
        - `links`
        - `responsible-roles`
          - `role-id`
          - `props`
          - `links`
          - `party-uuid`
          - `remarks`
        - `role-ids`
    - `by-components`
      - `satisfied`
        - `uuid`
        - `responsibility`
        - `description`
        - `props`
        - `links`
        - `responsible-roles`
          - `role-id`
          - `props`
          - `links`
          - `party-uuid`
          - `remarks`
        - `role-ids`

### Elements NOT in Component Definition that need added for SSP Generate

- `control-implementation`
  - `implemented-requirements`
    - `by-components`
      - `implementation-status`
        - `state`
        - `remarks`
    - `statements`
      - `by-components`
        - `implementation-status`
          - `state`
          - `remarks`

### Component Definition to SSP Transferable Fields

### **NOTE:** repetitive children elements have been truncated to reduce bloat

- `metadata`
  - `title`
  - `published`
  - `last-modified`
  - `version`
  - `oscal-version`
  - `revisions`
  - `document-ids`
  - `props`
  - `links`
  - `roles`
  - `locations`
  - `actions`
- `control-implementation`
  - `description`
  - `set-parameters`
    - `param-id`
    - `values`
    - `remarks`
  - `implemented-requirements`
    - `uuid`
    - `control-id`
    - `props`
    - `links`
    - `set-parameters`
      - `param-id`
      - `values`
      - `remarks`
    - `statements`
      - `statement-id`
      - `props`
      - `responsible-roles`
      - `links`
      - `by-components`
    - `by-components`
      - `component-uuid`
      - `description`
      - `props`
      - `links`
      - `set-parameters`
        - `param-id`
        - `values`
        - `remarks`
  - `system-implementation` (Contains Fields from Component Definition)
    - `components`
      - `uuid`
      - `type`
      - `title`
      - `description`
      - `purpose`
      - `props`
      - `links`
      - `status`
      - `protocols`
    - `implemented-components`
      - `component-uuid`
      - `props`
      - `links`
