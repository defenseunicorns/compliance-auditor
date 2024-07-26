# Validate

The `validate` command is used to run a single validation for the input Lula Validation manifest. This helps while developing validations to exercise the policy engines to ensure correctness and expose the outputs that will appear in the `assessment-results`.

## Usage

```bash
lula dev validate -f /path/to/validation.yaml
```

## Options

- `-f, --input-file`: The path to the target validation manifest.
- `-o, --output-file`: [Optional] The path to the output file. If not specified, the output will print to STDOUT
- `-r, --resources-file`: [Optional] The path to the resources file; must be json. If not specified, resources will be read from the domain.
- `-e, --expected-result`: [Optional] The expected result of the validation, true or false. Default is true.
- `-t, --timeout`: [Optional] Timeout when waiting for results from STDIN.
- `--confirm-execution`: [Optional] Flag to skip execution confirmation prompt. Only relevant when running a domain that performs some execution.

## Examples

To run validation using a custom resources file
```bash
lula dev validate -f /path/to/validation.yaml -r /path/to/resources.json
```

To run validation and automatically confirm execution
```bash
lula dev validate -f /path/to/validation.yaml --confirm-execution
```

To run validation from stdin
```bash
cat /path/to/validation.yaml | lula dev validate
```

To hang indefinitely for stdin
```bash
lula dev validate -t -1
```

To hang for timeout of 5 seconds
```bash
lula dev validate -t 5
```

