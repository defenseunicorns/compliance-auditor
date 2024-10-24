---
title: lula dev print-resources
description: Lula CLI command reference for <code>lula dev print-resources</code>.
type: docs
---
## lula dev print-resources

Print Resources from a Lula Validation evaluation

### Synopsis


Print out the the JSON resources input that were provided to a Lula Validation, as identified by a given observation and assessment results file.


```
lula dev print-resources [flags]
```

### Examples

```

To print resources from lula validation manifest:
	lula dev print-resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

To print resources from lula validation manifest to output file:
	lula dev print-resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --output-file /path/to/output.json

```

### Options

```
  -a, --assessment string         the path to an assessment-results file
  -h, --help                      help for print-resources
  -u, --observation-uuid string   the observation uuid
  -o, --output-file string        the path to write the resources json
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula dev](./lula_dev.md)	 - Collection of dev commands to make dev life easier

