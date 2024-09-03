---
title: lula tools lint
description: Lula CLI command reference for <code>lula tools lint</code>.
type: docs
---
## lula tools lint

Validate OSCAL against schema

### Synopsis

Validate OSCAL documents are properly configured against the OSCAL schema

```
lula tools lint [flags]
```

### Examples

```

To lint existing OSCAL files:
	lula tools lint -f <path1>,<path2>,<path3> [-r <result-file>]


```

### Options

```
  -h, --help                  help for lint
  -f, --input-files strings   the paths to oscal json schema files (comma-separated)
  -r, --result-file string    the path to write the validation result
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](/cli/cli-commands/lula_tools/)	 - Collection of additional commands to make OSCAL easier

