---
title: lula tools compose
description: Lula CLI command reference for <code>lula tools compose</code>.
type: docs
---
## lula tools compose

compose an OSCAL component definition

### Synopsis

Lula Composition of an OSCAL component definition. Used to compose remote validations within a component definition in order to resolve any references for portability.

```
lula tools compose [flags]
```

### Examples

```

To compose an OSCAL Model:
	lula tools compose -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools compose -f ./oscal-component.yaml -o composed-oscal-component.yaml

```

### Options

```
  -h, --help                    help for compose
  -f, --input-file string       the path to the target OSCAL component definition
  -o, --output-file -composed   the path to the output file. If not specified, the output file will be the original filename with -composed appended
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](/cli/cli-commands/lula_tools/)	 - Collection of additional commands to make OSCAL easier

