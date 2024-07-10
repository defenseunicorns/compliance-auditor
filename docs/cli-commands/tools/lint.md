# Lint Command

The `lint` command is used to validate OSCAL files against the OSCAL schema. It can validate both composed and non-composed OSCAL models.

## Usage

```bash
lula tools lint -f <input-files> [-r <result-file>] [-c]
```

## Options

- `-f, --input-files`: The paths to the target OSCAL files (comma-separated).
- `-r, --result-file`: The path to the result file. If not specified, the validation results will be printed to the console.
- `-c, --composed`: Disable composition before linting. Use this option only if you are sure that the OSCAL model is already composed (i.e., it has no imports or remote validations, default is false).

## Examples

To lint existing OSCAL files:
```bash
lula tools lint -f ./oscal-component1.yaml,./oscal-component2.yaml
```

To lint composed OSCAL models:
```bash
lula tools lint -c -f ./oscal-component1.yaml,./oscal-component2.yaml
```

To specify a result file:
```bash
lula tools lint -f ./oscal-component1.yaml,./oscal-component2.yaml -r validation-results.json
```

## Notes

If no input files are specified, an error will be returned. The validation results will be written to the specified result file. If no result file is specified, the validation results will be printed to the console. If there is at least one validation result that is not valid, the command will exit with a fatal error listing the files that failed linting.
