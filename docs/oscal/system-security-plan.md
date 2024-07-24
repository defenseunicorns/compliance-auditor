# System Security Plan

A [System Security Plan](https://pages.nist.gov/OSCAL-Reference/models/v1.1.2/system-security-plan/json-reference/#/system-security-plan) is an OSCAL-specific model to represent a system as a whole. In Lula, the `generate system-security-plan` command creates an `oscal-system-security-plan` object to explain the system as a whole by using the compliance data provided by the `component-definition` and `LulaOscalConfig`.

## Metadata

Includes all `responsible parties`, `parties`, and `roles` that plan a part in the system. Responsible parities are the parties who are responsible for the maintenance and development of the system. Parties includes any internal or external party that contribute to the system or the lifecycle of the system. Roles are the designated positions contributors take to within the system and system's lifecycle.

## System Characteristics

Describes the system and the systems security requirements. This includes the `security-sensitivity-level` which is the overall system's sensitivity categorization as defined by [FIPS-199](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.199.pdf). The system's overall level of expected impact resulting from unauthorized disclosure, modification, or loss of access to information through `security-impact-level` children items of `security-objective-confidentiality`, `security-objective-integrity`, and `security-objective-availability`.

The system characteristics also includes the `authorization-boundary`, `network-architecture`, and `data-flow` diagrams or links to the location of the diagrams used to describe the system.

## System Implementation

Contains any `leveraged-authorizations`, if used, all `components` used to build the system, all `users` with their type and access levels listed, and `inventory-items` detailing how the overall system is configured.

## Control Implementation

Contains all of the compliance controls the system must adhere to as outlined within the `profile`. Each `implemented-requirement` is listed detailing the control and the information of how the system meets the control on a `by-component` instance
