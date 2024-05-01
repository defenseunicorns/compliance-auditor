# Compose Command

The `compose` command is used to compose an OSCAL component definition. It is used to compose remote validations within a component definition in order to resolve any references for portability.

## Usage

```bash
lula tools compose -f <input-file> -o <output-file>
```

## Options

- `-f, --input-file`: The path to the target OSCAL component definition.
- `-o, --output-file`: The path to the output file. If not specified, the output file will be the original filename with `-composed` appended.

## Examples

To compose an OSCAL Model:
```bash
lula tools compose -f ./oscal-component.yaml
```

To indicate a specific output file:
```bash
lula tools compose -f ./oscal-component.yaml -o composed-oscal-component.yaml
```

## Notes

If the input file does not exist, an error will be returned. The composed OSCAL Component Definition will be written to the specified output file. If no output file is specified, the composed definition will be written to a file with the original filename and `-composed` appended.
