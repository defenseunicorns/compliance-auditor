# Contributing to Lula

Thank you for your interest in contributing to Lula! We welcome all contributions and are grateful for your help. This guide outlines how to get started with contributing to this project.

## Table of Contents

- [Contributing to Lula](#contributing-to-lula)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [Getting Started](#getting-started)
    - [Setup](#setup)
  - [Submitting a Pull Request](#submitting-a-pull-request)
    - [PR Requirements](#pr-requirements)
  - [Running Tests](#running-tests)
    - [Automated Tests](#automated-tests)
    - [Run Tests Locally](#run-tests-locally)
    - [Test a Local Development Version](#test-a-local-development-version)
  - [Documentation](#documentation)
    - [Architecture Design Records (ADR)](#architecture-design-records-adr)
    - [How to use `adr-tools`](#how-to-use-adr-tools)
  - [Contact](#contact)

## Code of Conduct

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md) to maintain a respectful and collaborative environment.

## Getting Started

- **Repository**: [https://github.com/defenseunicorns/lula/](https://github.com/defenseunicorns/lula/)
- **Go Binaries**: [https://github.com/defenseunicorns/lula/releases](https://github.com/defenseunicorns/lula/releases)
- **Required Go version**: `>=1.22.0`

### Setup

1. Fork the repository.
2. Clone your fork locally: `git clone https://github.com/your-username/lula.git`.
3. Create a new branch for your feature or fix: `git checkout -b my-feature-branch`.

## Submitting a Pull Request

1. **Create an Issue**: For significant changes, please create an issue first, describing the problem or feature proposal. Trivial fixes do not require an issue.
2. **Commit Your Changes**: Make your changes and commit them. All commits must be signed. For help follow this [guide.](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)
3. **Run Tests**: Ensure that your changes pass all tests by running `go test` for each required test.
4. **Push Your Branch**: Push your branch to your fork on GitHub.
5. **Create a Pull Request**: Open a pull request against the `main` branch of the lula repository. Please make sure that your PR passes all CI checks.

### PR Requirements

- PRs must be against the `main` branch.
- PRs must pass CI checks.
- All commits must be signed.
- PRs should have a related issue, except for trivial fixes.

## Running Tests

### Automated Tests

Automated tests will begin based on the paths you have edited in your Pull Request.
⚠️ NOTE: If you are an external third-party contributor, the pipelines won't run until a [CODEOWNER](CODEOWNERS) approves the pipeline run.

### Run Tests Locally

- Run E2E tests: `make test`

### Test a Local Development Version

1. Run `make build` and wait for completion.
2. Change to the lula directory: `cd bin/lula`.
3. You can now run any of the `lula` commands.

## Documentation

The decision for processes and docs on how-to can be found in `adr/` or in `docs/`.

### Architecture Design Records (ADR)

We've chosen to use ADRs to document architecturally significant decisions. We primarily use the guidance found in [this article by Michael Nygard](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions) with a couple of tweaks:

- The criteria for when an ADR is needed is undefined. The team will decide when the team needs an ADR.
- We will use the tool [adr-tools](https://github.com/npryce/adr-tools) to make it easier on us to create and maintain ADRs.
- We will keep ADRs in the repository under `adr/NNNN-name-of-adr.md`. `adr-tools` is configured with a dotfile to automatically use this directory and format.

### How to use `adr-tools`

```bash
# Create a new ADR titled "Use Bisquick for all waffle making"
adr new Use Bisquick for all waffle making

# Create a new ADR that supersedes a previous one. Let's say, for example, that the previous ADR about Bisquick was ADR number 9.
adr new -s 9 Use scratch ingredients for all waffle making

# Create a new ADR that amends a previous one. Let's say the previous one was ADR number 15
adr new -l "15:Amends:Amended by" Use store-bought butter for all waffle making

# Get full help docs. There are all sorts of other helpful commands that help manage the decision log.
adr help
```

## Contact

For any questions or concerns, please open an issue on GitHub or contact the maintainers.
