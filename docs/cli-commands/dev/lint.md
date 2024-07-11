# Lint Command

The `lula dev lint` command is used to validate validation files against the schema. It can validate both local files and URLs.

## Usage

```bash
lula dev lint -f <input-files> [-r <result-file>]
```

## Options

- `-f, --input-files`: The paths to the validation files (comma-separated).
- `-r, --result-file`: The path to the result file. If not specified, the validation results will be printed to the console.

## Examples

To lint existing validation files:
```bash
lula dev lint -f ./validation-file1.yaml,./validation-file2.yaml,https://example.com/validation-file3.yaml
```

To specify a result file:
```bash
lula dev lint -f ./validation-file1.yaml,./validation-file2.yaml -r validation-results.json
```

## Notes

The validation results will be written to the specified result file. If there is at least one validation result that is not valid, the command will exit with a fatal error listing the files that failed linting.
