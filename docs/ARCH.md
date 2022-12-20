# Architecture Document

The purpose of this document is to serve as a starting point to iterating on the architecture for Lula. Providing insight into requirements and preferred capabilities. 

## High Level Intent
As we began to prove out the automation workflows from OSCAL component definitions to declarative deployments, it was apparent that there were two models that are immediately applicable to compliance workflows. 

1. How can I determine if my configuration is compliant _before_ I introduce it to my runtime?
2. How can I continuously evaluate compliance of my live-environment?

### Control Specification

Optimally the specification for the controls that an application satisfies from a given standard or benchmark is both produced by/ lives alongside the project in the upstream. Maintainers and contributors to a project most often have the best understanding of the application runtime and would therefore understand the ability to satisfy a given control without the bounds of some context and configuration.

This creates the inheritance model where end-users can aggregate the documents for the applications that comprise their environment and derive any subset of controls satisfied (and not satisfied) given a configuration. 

### Execution

Given the above value statement - it drives the idea of two primary execution models:
- Command Line Interface
    - The tool that would allow for end-users to use in standard processes (development etc) as well as part of other automation such as CI/CD.
    - Innate ability to "audit" using reproducible and human-verifiable processes and tooling.
    - Static manifest 

- Kubernetes Operator/Controller
    - Providing the innate ability to enable **Continuous Compliance** for live-environments.
    - Validating configurations when a change is introduced and producing reports for Governance, Risk, and Compliance (GRC) tooling ingestion.
    - Providing gRPC (or other) endpoint for integration from other applications stacks that would benefit from direct-Lula integration.

### Active Enforcement

Lula should optimally not be an admission controller given the overlap between compliance and policy. Rather Lula should either provide some dynamic ability to create policies OR be positioned to integrate with policy engines in such a way to overlay control satisfaction into an enforceable policy should enforcement be preferred.

### Provenance

All environments will incur specific scenarios where there may be instances of exclusions that are required **AND** accepted as such by those evaluating the security posture of their system (and ultimately driving the use of Lula on their own system). This is why policy engines provide an **exclusions** interface for any given policy.

Given that understanding, one scenario that can adapt to this requirement is the ability for some position of authority to establish and sign for when a given exclusion is authorized. Actual implementation here may drive some complexity around what that actually looks like - but I would expect that Lula can explicitly require some signature for all controls - and any validation without some signature present would fail validation.

### OSCAL Processing

OSCAL is not a trivial format to interact with. There are no actively supported golang "libraries" in order to abstract the utilization for use with business logic. With that said, it does provide a standard that can be worked against in order to provide a common definition for use internally and externally. CNCF groups are evaluating the standardization of control information in sponsored projects through the use of tooling and OSCAL documentation. 

#### Short-Term / low-oscal-integration

If Lula were to simply ingest OSCAL component definitions, execute validation, and produce assessment results. This could ease the ability to maintain OSCAL datatypes across version updates through generation of datatypes from the OSCAL json schema.

See special note for integration with Visualization below.

#### Long-Term / high-oscal-integration

Providing Golang libraries for both Metaschema and OSCAL would be an enabler for both the OSCAL ecosystem, as well as tooling such as Lula. If Lula is to ingest/produce any large subsets of the OSCAL models (Catalog/SSP/Component Definition/Assessment Plan/Assessment Results/ etc) across different versions of OSCAL documents over an extended lifecycle - then it may behoove the tooling to have underlying support through a [metaschema](https://pages.nist.gov/metaschema/specification/) implementation as opposed to manually generated structs.

### Other Processing

Support could be provided for ingestion of multiple formats to include OSCAL and other more-developer-friendly formats. Given the thoughts around integration, it seems reasonable that the tooling could ingest other data formats that expressed the relationship between expected configuration and actual configuration. 

### Visualization Integration

Definition of a standard approach for exposing data for visualization will be required. These could be:
- A plugin for coordinating data upload
- An external "shipper" for coordinating upload outside the boundary
- Expectation that the visualization component is going to call Lula through an endpoint for the data

* Note: The ability to leverage a visualization tool that understands OSCAL models could remove the need for Lula to have deep integration with OSCAL. 