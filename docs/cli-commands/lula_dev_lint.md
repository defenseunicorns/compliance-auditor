---
title: lula dev lint
description: Lula CLI command reference for <code>lula dev lint</code>.
type: docs
---
## lula dev lint

Lint validation files against schema

### Synopsis

Validate validation files are properly configured against the schema, file paths can be local or URLs (https://)

```
lula dev lint [flags]
```

### Examples

```

To lint existing validation files:
	lula dev lint -f <path1>,<path2>,<path3> [-r <result-file>]

```

### Options

```
  -h, --help                  help for lint
  -f, --input-files strings   the paths to validation files (comma-separated)
  -r, --result-file string    the path to write the validation result
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
  -s, --set strings        set a value in the template data
```

### SEE ALSO

* [lula dev](./lula_dev.md)	 - Collection of dev commands to make dev life easier

