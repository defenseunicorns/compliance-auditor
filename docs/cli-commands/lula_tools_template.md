---
title: lula tools template
description: Lula CLI command reference for <code>lula tools template</code>.
type: docs
---
## lula tools template

Template an artifact

### Synopsis

Resolving templated artifacts with configuration data

```
lula tools template [flags]
```

### Examples

```

To template an OSCAL Model:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

Data for the templating should be stored under the 'variables' configuration item in a lula-config.yaml file

```

### Options

```
  -h, --help                 help for template
  -f, --input-file string    the path to the target artifact
  -o, --output-file string   the path to the output file. If not specified, the output file will be directed to stdout
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](./lula_tools.md)	 - Collection of additional commands to make OSCAL easier

