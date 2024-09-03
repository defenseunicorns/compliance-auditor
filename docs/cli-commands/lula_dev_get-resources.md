---
title: lula dev get-resources
description: Lula CLI command reference for <code>lula dev get-resources</code>.
type: docs
---
## lula dev get-resources

Get Resources from a Lula Validation Manifest

### Synopsis

Get the JSON resources specified in a Lula Validation Manifest

```
lula dev get-resources [flags]
```

### Examples

```

To get resources from lula validation manifest:
	lula dev get-resources -f /path/to/validation.yaml
To get resources from lula validation manifest and write to file:
	lula dev get-resources -f /path/to/validation.yaml -o /path/to/output.json
To get resources from lula validation and automatically confirm execution
	lula dev get-resources -f /path/to/validation.yaml --confirm-execution
To run validations using stdin:
	cat /path/to/validation.yaml | lula dev get-resources
To hang indefinitely for stdin:
	lula get-resources -t -1
To hang for timeout of 5 seconds:
	lula get-resources -t 5

```

### Options

```
      --confirm-execution    confirm execution scripts run as part of getting resources
  -h, --help                 help for get-resources
  -f, --input-file string    the path to a validation manifest file (default "0")
  -o, --output-file string   the path to write the resources json
  -t, --timeout int          the timeout for stdin (in seconds, -1 for no timeout) (default 1)
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula dev](/cli/cli-commands/lula_dev/)	 - Collection of dev commands to make dev life easier

