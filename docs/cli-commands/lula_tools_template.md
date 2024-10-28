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

To template an OSCAL Model, defaults to masking sensitive variables:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

To perform overrides on the template data:
	lula tools template -f ./oscal-component.yaml --set .var.key1=value1 --set .const.key2=value2

To perform the full template operation, including sensitive data:
	lula tools template -f ./oscal-component.yaml --render all

Data for templating should be stored under 'constants' or 'variables' configuration items in a lula-config.yaml file
See documentation for more detail on configuration schema

```

### Options

```
  -h, --help                 help for template
  -f, --input-file string    the path to the target artifact
  -o, --output-file string   the path to the output file. If not specified, the output file will be directed to stdout
  -r, --render string        values to render the template with, options are: masked, constants, non-sensitive, all (default "masked")
  -s, --set strings          set a value in the template data
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](./lula_tools.md)	 - Collection of additional commands to make OSCAL easier

