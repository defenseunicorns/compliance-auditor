---
title: lula dev print-validation
description: Lula CLI command reference for <code>lula dev print-validation</code>.
type: docs
---
## lula dev print-validation

Print Lula Validation

### Synopsis


Prints the Lula Validation from a specified observation. Assumes that the validation is in the back matter of the provided component definition.


```
lula dev print-validation [flags]
```

### Examples

```

To print a specific lula validation that generated a given observation:
	lula dev print-validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

```

### Options

```
  -a, --assessment string         the path to an assessment-results file
  -c, --component string          the path to a validation manifest file
  -h, --help                      help for print-validation
  -u, --observation-uuid string   the observation uuid
  -o, --output-file string        the path to write the resources json
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula dev](./lula_dev.md)	 - Collection of dev commands to make dev life easier

