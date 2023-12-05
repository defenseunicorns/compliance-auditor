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
    - Kind
        - `kind create cluster -n lula-test`
    - K3d
        - `k3d cluster create lula-test`
- kubectl
- GoLang version 1.21.x

#### Steps

1. Clone the repository to your local machine and change into the `lula` directory

    ```shell
    git clone https://github.com/defenseunicorns/lula.git && cd lula
    ```

2. While in the `lula` directory, compile the tool into an executable binary. This outputs the `lula` binary to the `bin` directory.

    ```shell
    make build
    ```

3. Apply the `./demo/namespace.yaml` file to create a namespace for the demo

    ```shell
    kubectl apply -f ./demo/namespace.yaml
    ```

4. Apply the `./demo/pod.fail.yaml` to create a pod in your cluster

    ```shell
    kubectl apply -f ./demo/pod.fail.yaml
    ```

5. Run the following command in the `lula` directory:

    ```shell
    ./bin/lula validate ./demo/oscal-component.yaml
    ```

    The output in your terminal should inform you that the control validated is `not-satisfied`:

    ```shell
    OPA provider validating...
    UUID: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
            Status: not-satisfied
    ```

    This will also produce an assessment-results file with timestamp - review the findings and observations:

    ```yaml
    findings:
    - description: Lorem ipsum dolor sit amet, consectetur adipiscing elit....
      related-observations:
        - observation-uuid: 51fe298d-16b9-4efb-9a0f-f3ab54da50af
      target:
        status:
            state: not-satisfied
        target-id: ID-1
        type: objective-id
      title: 'Validation Result - Component:A9D5204C-7E5B-4C43-BD49-34DF759B9F04 / Control Implementation: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A / Control:  ID-1'
      uuid: 32ad2bce-e2f6-4445-a96e-a3b693b942f1
    observations:
    - collected: "2023-12-01T13:22:09-08:00"
      description: |
        [TEST] ID-1 - a7377430-2328-4dc4-a9e2-b3f31dc1dff9
      methods:
        - TEST
      relevant-evidence:
        - description: |
            Result: not-satisfied - Passing Resources: 0 - Failing Resources 1
      uuid: 51fe298d-16b9-4efb-9a0f-f3ab54da50af
    ```

6. Now, apply the `./demo/pod.pass.yaml` file to your cluster to configure the pod to pass compliance validation:

    ```shell
    kubectl apply -f ./demo/pod.pass.yaml
    ```

7. Run the following command in the `lula` directory:

    ```shell
    ./bin/lula validate ./demo/oscal-component.yaml
    ```

    The output should now show the pod as passing the compliance requirement:

    ```shell
    OPA provider validating...
    UUID: 42C2FFDC-5F05-44DF-A67F-EEC8660AEFFD
            Status: satisfied
    ```

    This will produce a new assessment-results file with timestamp - review the findings and observations:

    ```yaml
    findings:
    - description: Lorem ipsum dolor sit amet, consectetur adipiscing elit...
      related-observations:
        - observation-uuid: 51fe298d-16b9-4efb-9a0f-f3ab54da50af
      target:
        status:
            state: not-satisfied
        target-id: ID-1
        type: objective-id
      title: 'Validation Result - Component:A9D5204C-7E5B-4C43-BD49-34DF759B9F04 / Control Implementation: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A / Control:  ID-1'
      uuid: 32ad2bce-e2f6-4445-a96e-a3b693b942f1
    observations:
    - collected: "2023-12-01T13:22:09-08:00"
      description: |
        [TEST] ID-1 - a7377430-2328-4dc4-a9e2-b3f31dc1dff9
      methods:
        - TEST
      relevant-evidence:
        - description: |
            Result: not-satisfied - Passing Resources: 0 - Failing Resources 1
      uuid: 51fe298d-16b9-4efb-9a0f-f3ab54da50af
    ```

## Future Extensibility

- Support for cloud infrastructure state queries

## Developing

- Go 1.21
