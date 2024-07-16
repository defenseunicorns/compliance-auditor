# Lula in CI

Lula is designed to evaluate the _continual_ compliance of a system, and as such is a valuable tool to implement in a CI environment to provide rapid feedback to developers if a system moves out of compliance. The [Lula-Action Repo](https://github.com/defenseunicorns/lula-action) supports the use of Lula in github workflows and this document provides an outline for implementation.

### Pre-Requisite

To use Lula to `validate` and `evaluate` a system _in development_, a pre-requisite is having an OSCAL Component Definition model, along with linked `Lula Validations` existing in the repository, a sample structure follows:

```bash
.
|-- .github
|   |-- workflows
|-- |-- |-- lint.yaml # Existing workflow to lint
|-- |-- |-- test.yaml # Existing workflow to test system
|-- README.md
|-- LICENSE
|-- compliance
|-- |-- oscal-component.yaml # OSCAL Component Definition
|-- src
|   |-- main
|   |-- test
```

### Steps

1) Add Lula linting to `.github/workflows/lint.yaml`:

    ```yaml
    name: Lint

    on:
        pull_request:
            branches: [ "main" ]

    jobs:
        lint:
            runs-on: ubuntu-latest

            # ... Other jobs

            - name: Setup Lula
              uses: defenseunicorns/lula-action/setup@main
              with:
                version: v0.4.1
            
            - name: Lint OSCAL file
              uses: defenseunicorns/lula-action/lint@main
              with:
                oscal-target: ./compliance/oscal-component.yaml
            
            # ... Other jobs
    ```

    Additional linting targets may be added to this list as comma separated values, e.g., `component1.yaml,component2.yaml`. Note that linting is only validating the correctness of the OSCAL.

2) Add Lula validation and evaluation to the testing workflow, `.github/workflows/test.yaml`:

    ```yaml
    name: Test

    on:
    pull_request:
        branches: [ "main" ]

    jobs:
        test:
            runs-on: ubuntu-latest

            # ... Other jobs

            - name: Setup Lula
              uses: defenseunicorns/lula-action/setup@main
              with:
                version: v0.4.1
            
            - name: Validate Component Definition
              uses: defenseunicorns/lula-action/validate@main
              with:
                oscal-target: ./compliance/oscal-component.yaml
                threshold: ./assessment-results.yaml
            
            # ... Other jobs
        test-upgrade:
            runs-on: ubuntu-latest

            # ... Jobs to deploy previous system version 

            - name: Setup Lula
              uses: defenseunicorns/lula-action/setup@main
              with:
                version: v0.4.1
            
            - name: Validate Component Definition
              uses: defenseunicorns/lula-action/validate@main
              with:
                oscal-target: ./compliance/oscal-component.yaml
                threshold: ./assessment-results.yaml
            
            # ... Jobs to upgrade system to current version
    ```

    The first `validate` under `test` outputs an `assessment-results` model that provide the assessment of the system in the current state. The second `validate` that occurs in the `test-upgrade` job runs a validation on the previous version of the system prior to upgrade. It then compares the old and new assessment results to either pass or fail the job - failure occurs when the current system's compliance is **worse** than the old system.