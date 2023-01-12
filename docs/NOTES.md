# Notes

## Foundation
- Golang OSCAL support for models/types
    - OSCAL developers recommend generation of the Go Types from metaschema as the most maintainable solution
        - [Example](https://github.com/GoComply/metaschema)
        - There will be a learning curve involved here. Not insurmountable, but will require time.
    - Second option is to follow other tooling such as [Compliance-Trestle](https://github.com/IBM/compliance-trestle)
        - Generate model types from a released NIST JSON Schema
        - This may provide a good first iteration with the smallest amount of work required to get usable types
            - un-blocks implementations that would benefit from consuming this as a module
        - Should be located in an independent repository for maintenance purposes
            - Would provide a iterative location should the work be converted to use metaschema instead of OSCAL JSON schema
    - Type Generation in isolation to reduce dependency requirements
        - A separate CLI that generates oscal type file(s) for use in Lula - meaning that Lula would not require any modules that perform the generation steps - as this would not be an action that Lula would perform. 

- Understanding (and documenting) the full end-to-end intent for compliance
    - The relationship between components and platforms (and how an end-user would leverage component models for an SSP)
    - The value that a single application (optimally open source) could benefit by collaborating/producing OSCAL
    - The value to end-users for inheritable controls with continuous-compliance validation
    - the value to GRC tooling for integration (see more below)

- Understanding the ecosystem
    - How might tools want to integrate with Lula for leveraging functionality in application-specific use cases
        - gRPC being an integration path for validation without orchestration
    - There are spaces with which Lula should not look to compete - IE policy enforcement engines
        - But there is a lot of alignment opportunity - how?
    - Plugin architecture
        - May cross with the mention of gRPC above
        - Can we optimize to allow the data to get where it needs
            - GRC Tooling will allow relationships to be established between models (if present)
            - Other tools might want to send an OSCAL (or other) payload to be validated with results returned for some catered user experience.

- How to validate configurations
    - We know kubernetes orchestration is declarative-in-nature and that expected configuration can be validated against actual configuration via the API-Server
    - Overlap between compliance and policy is quite heavy - leveraging policy to perform compliance-validation potentially improves both
        - Leveraging a policy engine to perform configuration queries provides an interfacing language (rego etc) to utilize in OSCAL validation rules

- Lula Run Time    
    - The creation of a Command Line Interface (CLI) for the purpose of validating static manifests that are not currently deployed to a live-environment.
    - A controller/operator that can continuously validate and monitor for changes in resources to report on changes in compliance.

- How is Application Compliance discoverable
    - OSCAL documentation for software components should live with the codebase
        _ Should this be application source or orchestration (Which can be separate)
    - How will Lula discover these compliance documents?
        - OSCAL is a common format that we can utilize 
        - A CRD is straight-forward, but integration would require that CRD be declared on the runtime environment - which is not always true
            - A CRD will be tough to persuade Open Source Projects to adopt
        - It would be worth-while to build out some examples that provide exposing compliance information via REST/gRPC from the perspective of an application owner.
        - Kubernetes Native resources may provide one solution for deployment of the compliance data with the application. In much the same way that grafana can discover dashboards by configmap labels as seen [Here](https://johnharris.io/2019/03/dynamic-configuration-discovery-in-grafana/)
            - This is still constrained to a single namespace - Lula will likely have the ability/permissions to see cluster globally - so maybe the resource can simply be deployed in application target namespace and Lula can discover across all namespaces.

- Non-Kubernetes compliance validation to perform
    - Infrastructure as Code (IaC)
        - We can provide Lula with a runtime environment on Kubernetes - it can then be extended to perform validation of the Infrastructure that it has connectivity-to.
        - There will be complexity in defining what validates or invalidates a control that infrastructure provides. 
        - IaC will use a component-definition OSCAL much like Kubernetes
    - API/OS/etc
        - More thoughts to come on other domains of controls that can be captured and continuously validated.

## Extension
- Existence of the tool/resource(s)
    - Can we fix this issue in the Kyverno CLI?
    - Or do we provide a different layer for providing this validation?

## Kyverno

### Limitations
- Wildcard for match any "kind" does not work as specified
- Cannot correlate between multiple resources well
    - Example: Given global context - cannot verify that all namespaces include a particular resource (CRD existence may be important)
    - Review the use of [External Data Sources](https://kyverno.io/docs/writing-policies/external-data-sources/) through the CLI
        - There are known issues here - but if resolvable could provide a mechanism for reporting on the existence of a resource
        - This would also allow establishing 1 -> Many relationships among resources