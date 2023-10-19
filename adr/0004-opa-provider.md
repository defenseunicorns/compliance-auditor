# 4. opa-provider

Date: 2023-09-26

## Status

Accepted

## Context

[Open Policy Agent](https://www.openpolicyagent.org/docs/latest/) (OPA) provides a general-purpose policy engine for use across various domains, including Kubernetes, infrastructure, APIs, and more. One of OPA's standout features is its ability to accept arbitrary structured data as input, primarily in JSON format. This capability, leveraged through the OPA SDK, empowers us to create data sources that function as queries, generating structured data, and leverage a general-purpose language for performing validations.

## Decision

The decision is to integrate Open Policy Agent (OPA) as a Provider through use of the OPA SDK within Lula. This integration will enable Lula to leverage OPA's capabilities in validating structured data, thereby enhancing its functionality in ensuring compliance and adherence to policies within environments.

## Consequences

The introduction of OPA as a validator brings several significant consequences and benefits:

- **Versatile Validation:** Leveraging OPA's capabilities allows for versatile and context-aware validation of resources. This flexibility empowers Lula to enforce policies and constraints effectively.

- **Structured Data Input:** OPA's ability to accept arbitrary structured data as input aligns perfectly with Lula's requirements for data collection and separate data collection functionality.

- **Interoperability:** By integrating OPA, Lula establishes a baseline for integrating data sources and validator functionality. This enables developers to create validation payloads that encompass various aspects of data source provisioning and their relationships with other domains, such as cloud infrastructure.

- **Scalability:** OPA's general-purpose language and versatility make it suitable for handling complex validation scenarios, making Lula more scalable in terms of policy enforcement and compliance checks.

- **Learning Curve:** It's important to acknowledge that integrating OPA introduces a learning curve for Lula developers who are not familiar with OPA's policy language. However, this learning investment can yield substantial benefits in terms of validation capabilities.
