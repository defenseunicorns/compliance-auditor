# About

Lula is a tool designed to bridge the gap between expected configuration required for compliance and **_actual_** configuration.

### Key Features
* **Assess** compliance of a system against user-defined controls
* **Evaluate** an evolving system for compliance _over time_
* **Generate** machine-readible OSCAL artifacts
* **Accelerate** the compliance and accreditation process

### Why Lula is different than a standard policy engine
* Lula is not meant to compete with policy engines - rather augment the auditing and alerting process
* Often admission control processes have a difficult time establishing `big picture` global context control satisfaction, Lula fills this gap
* Lula is meant to allow modularity and inheritance of controls based upon the components of the system you build

## Overview

Cloud-Native Infrastructure, Platforms, and Applications can establish [OSCAL documents](https://pages.nist.gov/OSCAL/about/) that are maintained alongside source-of-truth code bases. These documents provide an inheritance model to prove when a control that the technology can satisfy _IS_ satisfied in a live-environment.

These controls can be well established and regulated standards such as NIST 800-53. They can also be best practices, Enterprise Standards, or simply team development standards that need to be continuously monitored and validated.

Lula operates on a framework of proof by adding custom overlays mapped to the these controls, [`Lula Validations`](link), to measure system compliance. These `Validations` are constructed by establishing the collection of measurements about a system, given by the specified **Domain**, and the evaluation of adherence, performed by the **Provider**. 

**Domain** is the identifier for where and which data to collect as "evidence". Below are the active and planned domains:

| Domain | Current | Roadmap |
|----------|----------|----------|
| [Kubernetes](./docs/reference/domains/kubernetes.md) | ✅ | - |
| [API](./docs/reference/domains/api-domain.md) | ✅ | - |
| Cloud Infrastructure | ❌ | ✅ |

**Provider** is the "engine" performing the validation using policy and the data collected. Below are the active providers:

| Provider | Current | Roadmap |
|----------|----------|----------|
| [OPA](./docs/reference/provideres/opa-provider.md) | ✅ | - |
| [Kyverno](./docs/reference/provideres/kyverno-provider.md) | ✅ | - |