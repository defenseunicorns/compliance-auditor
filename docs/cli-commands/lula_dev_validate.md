---
title: lula dev validate
description: Lula CLI command reference for <code>lula dev validate</code>.
type: docs
---
## lula dev validate

Run an individual Lula validation.

### Synopsis

Run an individual Lula validation for quick testing and debugging of a Lula Validation. This command is intended for development purposes only.

```
lula dev validate [flags]
```

### Examples

```

To run validation from a lula validation manifest:
	lula dev validate -f /path/to/validation.yaml
To run validation using a custom resources file:
	lula dev validate -f /path/to/validation.yaml -r /path/to/resources.json
To run validation and automatically confirm execution
	lula dev validate -f /path/to/validation.yaml --confirm-execution
To run validation from stdin:
	cat /path/to/validation.yaml | lula dev validate
To hang indefinitely for stdin:
	lula dev validate -t -1
To hang for timeout of 5 seconds:
	lula dev validate -t 5

```

### Options

```
      --confirm-execution       confirm execution scripts run as part of the validation
  -e, --expected-result         the expected result of the validation (-e=false for failing result) (default true)
  -h, --help                    help for validate
  -f, --input-file string       the path to a validation manifest file (default "0")
  -o, --output-file string      the path to write the validation with results
  -r, --resources-file string   the path to an optional resources file
  -t, --timeout int             the timeout for stdin (in seconds, -1 for no timeout) (default 1)
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula dev](/cli/cli-commands/lula_dev/)	 - Collection of dev commands to make dev life easier

