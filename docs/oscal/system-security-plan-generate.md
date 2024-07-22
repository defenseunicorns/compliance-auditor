# System-Security-Plan Generation

## Description

The `system-security-plan` can be generated using the upstream catalog/profile in conjunction with the `component-definition`. There are net new fields that are apart of the `system-security-plan` that are not within the `component-definition` or catalog/profile that do not make sense to add as props. Those items are under the section `Elements in SSP Not in Component Definition`. There are items that are not in the `system-security-plan` but also not in the `component-definition` that do make sense to create as props. Those items are under the section `Elements NOT in Component Definition that need added for SSP Generate`. Lastly as a note there are items within the `component-definition` that are not used in the `system-security-plan`.

For the items in `Elements in SSP Not in Component Definition` a "config" file will be needed to fill in the gaps. There are also metadata fields such as `responsible-roles`, `responsible-parties`, and `parties` that can be added to the `system-security-plan` through the config file that may not be necessary to add directly to the `component-definition`

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
