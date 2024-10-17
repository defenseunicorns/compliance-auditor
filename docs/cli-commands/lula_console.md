---
title: lula console
description: Lula CLI command reference for <code>lula console</code>.
type: docs
---
## lula console

Console terminal user interface for OSCAL models

### Synopsis


The Lula Console is a text-based terminal user interface that allows users to 
interact with the OSCAL documents in a more intuitive and visual way.


```
lula console [flags]
```

### Examples

```

To view an OSCAL model in the Console:
	lula console -f /path/to/oscal-component.yaml

To view multiple OSCAL models in the Console:
	lula console -f /path/to/oscal-component.yaml,/path/to/oscal-assessment-results.yaml

To specify an output file to save any changes made to the component definition:
	lula console -f /path/to/oscal-component.yaml -c /path/to/output.yaml

```

### Options

```
  -c, --component-output string   the path to the component definition output file
  -h, --help                      help for console
  -f, --input-files strings       the path to the target OSCAL models, comma separated
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula](./lula.md)	 - Risk Management as Code

