# Validate

The `validate` command aims to bridge the gap between fully developing the `assessment-plan` and allowing for testing of the existing `component-definition`.

## Usage

```bash
lula validate -f /path/to/oscal-component.yaml
```

## Options

- `-f, --input-file`: The path to the OSCAL component definition file.
- `-o, --output-file`: [Optional] The path to the output assessment results file. Creates a new file or appends to existing. If not specified, the output file is `./assessment-results.yaml`
- `--confirm-execution`: [Optional] Flag to skip execution confirmation prompt. Only relevant when running a validation with a domain that performs some execution.
- `--non-interactive`: [Optional] Flag to indicate running non-interactively, i.e., does not request user to confirm validations with execution.

## Examples

This command is used both locally as an evaluation of the Component Definition to understand the component's compliance. It's also implemented in CI workflows to continually evaluate the evolution of a system during development. See the following relevant tutorials:

- ...