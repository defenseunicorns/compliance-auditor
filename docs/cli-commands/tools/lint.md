# Lint Command

The `lula tools lint` command is used to validate OSCAL files against the OSCAL schema. It can validate both composed and non-composed OSCAL models.
> **Note**: the `lint` command does not compose the OSCAL model.
> If you want to validate a composed OSCAL model, you should use the [`lula tools compose`](./compose/README.md) command first.

## Usage

```bash
lula tools lint -f <input-files> [-r <result-file>]
```

## Options

- `-f, --input-files`: The paths to the tar get OSCAL files (comma-separated).
- `-r, --result-file`: The path to the result file. If not specified, the validation results will be printed to the console.

## Examples

To lint existing OSCAL files:
```bash
lula tools lint -f ./oscal-component1.yaml,./oscal-component2.yaml
```

To specify a result file:
```bash
lula tools lint -f ./oscal-component1.yaml,./oscal-component2.yaml -r validation-results.json
```

## Notes

If no input files are specified, an error will be returned. The validation results will be written to the specified result file. If no result file is specified, the validation results will be printed to the console. If there is at least one validation result that is not valid, the command will exit with a fatal error listing the files that failed linting.