# Get Resources

The `get-resources` command is used to execute the domain portion of the input Lula Validation manifest to extract the `resources` json. This helps while developing validations to ensure the domain specification is returning the expected resources.

## Usage

```bash
lula dev get-resources -f /path/to/validation.yaml
```

## Options

- `-f, --input-file`: The path to the target validation manifest.
- `-o, --output-file`: [Optional] The path to the output file. If not specified, the output will print to STDOUT
- `-t, --timeout`: [Optional] Timeout when waiting for results from STDIN.
- `--confirm-execution`: [Optional] Flag to skip execution confirmation prompt. Only relevant when running a domain that performs some execution.

## Examples

To get resources and write to file
```bash
lula dev get-resources -f /path/to/validation.yaml -o /path/to/resources.json
```

To run get resources and automatically confirm execution
```bash
lula dev get-resources -f /path/to/validation.yaml --confirm-execution
```

To run get resources from stdin
```bash
cat /path/to/validation.yaml | lula dev get-resources
```

To hang indefinitely for stdin
```bash
lula dev get-resources -t -1
```

To hang for timeout of 5 seconds
```bash
lula dev get-resources -t 5
```

