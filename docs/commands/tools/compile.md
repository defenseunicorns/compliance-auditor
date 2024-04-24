# Compile Command

The `compile` command is used to compile an OSCAL component definition. It is used to compile remote validations within a component definition in order to resolve any references for portability.

## Usage

```bash
lula tools compile -f <input-file> -o <output-file>
```

## Options

- `-f, --input-file`: The path to the target OSCAL component definition.
- `-o, --output-file`: The path to the output file. If not specified, the output file will be the original filename with `-compiled` appended.

## Examples

To compile an OSCAL Model:
```bash
lula tools compile -f ./oscal-component.yaml
```

To indicate a specific output file:
```bash
lula tools compile -f ./oscal-component.yaml -o compiled-oscal-component.yaml
```

## Notes

If the input file does not exist, an error will be returned. The compiled OSCAL Component Definition will be written to the specified output file. If no output file is specified, the compiled definition will be written to a file with the original filename and `-compiled` appended.
