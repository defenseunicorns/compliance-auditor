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

```

### Options

```
  -h, --help                help for console
  -f, --input-file string   the path to the target OSCAL model
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula](/cli/cli-commands/lula/)	 - Risk Management as Code

