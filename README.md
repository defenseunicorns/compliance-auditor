# Lula - The Kubernetes Compliance Engine

[![Go version](https://img.shields.io/github/go-mod/go-version/defenseunicorns/lula?filename=go.mod)](https://go.dev/)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/defenseunicorns/lula/badge)](https://api.securityscorecards.dev/projects/github.com/defenseunicorns/lula)

Lula is a tool written to bridge the gap between expected configuration required for compliance and **_actual_** configuration.

Cloud Native Infrastructure, Platforms, and applications can establish OSCAL documents that live beside source-of-truth code bases. Providing an inheritance model for when a control that the technology can satisfy _IS_ satisfied in a live-environment.

This can be well established and regulated standards such as NIST 800-53. It can also be best practices, Enterprise Standards, or simply team development standards that need to be continuously monitored and validated.

## Why this approach vs a policy engine?

- Lula is not meant to compete with policy engines - rather augment the auditing and alerting process
- Often admission control processes have a difficult time establishing `big picture` global context control satisfaction
- Lula is meant to allow modularity and inheritance of controls based upon the components of the system you build

## How does it work?

Under the hood, Lula has two primary capabilities; Provider and Domains.

- A Domain is an identifier for where to collect data to be validated
- A Provider is the "engine" performing the validation using policy and the data collected.

In the standard CLI workflow:

- Target a `Component-Definition` OSCAL file for validation
  - `lula validate oscal-component.yaml`
- This creates an object in memory for the OSCAL content
- Lula then traverses as required to identify `implemented-requirements` that contain a Lula Validation Payload
- When the payload has been identified:
  - Lula processes provider to understand which provider to use for validation
    - More than one provider can be used in an OSCAL document
  - Lula processes the domain to understand how data is collected (and which data to collect)
  - Lula collects the data for validation as specified in the payload
  - Lula performs validation of the data collected as specified as policy in the payload

## Getting Started

### Try it out

#### Dependencies

- A running Kubernetes cluster
- GoLang version 1.21.x

#### Steps

1. Clone the repository to your local machine and change into the `lula` directory

    ```shell
    git clone https://github.com/defenseunicorns/lula.git && cd lula
    ```

1. While in the `lula` directory, compile the tool into an executable binary. This outputs the `lula` binary to the `bin` directory.

    ```shell
    make build
    ```

1. Apply the `./demo/namespace.yaml` file to create a namespace for the demo

    ```shell
    kubectl apply -f ./demo/namespace.yaml
    ```

1. Apply the `./demo/pod.fail.yaml` to create a pod in your cluster

    ```shell
    kubectl apply -f ./demo/pod.fail.yaml
    ```

1. Run the following command in the `lula` directory:

    ```shell
    ./lula validate ./demo/oscal-component.yaml
    ```

    The output in your terminal should inform you that there is at least one failing pod in the cluster:

    ```shell
    Applying 1 policy rule to 19 resources...

    policy 42c2ffdc-5f05-44df-a67f-eec8660aeffd -> resource foo/Pod/demo-pod failed: 
    1. ID-1: validation error: Every pod in namespace 'foo' should have 'foo=bar' label. rule ID-1 failed at path /metadata/labels/foo/ 
    UUID: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
            Resources Passing: 0
            Resources Failing: 1
            Status: Fail
    ```

1. Now, apply the `./demo/pod.pass.yaml` file to your cluster to configure the pod to pass compliance validation:

    ```shell
    kubectl apply -f ./demo/pod.pass.yaml
    ```

1. Run the following command in the `lula` directory:

    ```shell
    ./lula validate ./demo/oscal-component.yaml
    ```

    The output should now show the pod as passing the compliance requirement:

    ```shell
    Applying 1 policy rule to 19 resources...
    UUID: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
            Resources Passing: 1
            Resources Failing: 0
            Status: Pass
    ```

## Future Extensibility

- Support for cloud infrastructure state queries
- Support for API validation

## Developing

- Go 1.21
