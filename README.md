# Lula - The Cloud-Native Compliance Engine

[![Lula Documentation](https://img.shields.io/badge/docs--d25ba1)](https://docs.lula.dev)
[![Go version](https://img.shields.io/github/go-mod/go-version/defenseunicorns/lula?filename=go.mod)](https://go.dev/)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/defenseunicorns/lula/badge)](https://api.securityscorecards.dev/projects/github.com/defenseunicorns/lula)

<img align="right" src="./images/lula.svg" alt="lula logo" style="width:25%; height:auto;">

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

Lula operates on a framework of proof by adding custom overlays mapped to the these controls, [`Lula Validations`](./docs/reference/README.md), to measure system compliance. These `Validations` are constructed by establishing the collection of measurements about a system, given by the specified **Domain**, and the evaluation of adherence, performed by the **Provider**. 

### Providers and Domains

**Domain** is the identifier for where and which data to collect as "evidence". Below are the active and planned domains:

| Domain | Current | Roadmap |
|----------|----------|----------|
| [Kubernetes](./docs/reference/domains/kubernetes-domain.md) | ✅ | - |
| [API](./docs/reference/domains/api-domain.md) | ✅ | - |
| [File](./docs/reference/domains/file-domain.md) | ✅ | - |
| Cloud Infrastructure | ❌ | ✅ |

**Provider** is the "engine" performing the validation using policy and the data collected. Below are the active providers:

| Provider | Current | Roadmap |
|----------|----------|----------|
| [OPA](./docs/reference/providers/opa-provider.md) | ✅ | - |
| [Kyverno](./docs/reference/providers/kyverno-provider.md) | ✅ | - |

## Getting Started

[Install Lula](./docs/getting-started/README.md) and check out the [Simple Demo](./docs/getting-started/simple-demo.md) to get familiar with Lula's `validate` and `evaluate` workflow to assess system compliance and establish thresholds. See the other tutorials for more advanced Lula use cases and information on how to develop your own `Lula Validations`! 

## Communication

For more information on how to get involved in the community, mailing lists and
meetings, please refer to our [community page](./docs/community-and-contribution/README.md)

For security issues or code of conduct concerns, an e-mail should be sent to
[lula@defenseunicorns.com](mailto:lula@defenseunicorns.com).