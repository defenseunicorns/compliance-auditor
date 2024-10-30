---
title: lula dev print
description: Lula CLI command reference for <code>lula dev print</code>.
type: docs
---
## lula dev print

Print Resources or Lula Validation from an Assessment Observation

### Synopsis


Print out data about an Observation. 
Given "--resources", the command will print the JSON resources input that were provided to a Lula Validation, as identified by a given observation and assessment results file. 
Given "--validation", the command will print the Lula Validation that generated a given observation, as identified by a given observation, assessment results file, and component definition file.


```
lula dev print [flags]
```

### Examples

```

To print resources from lula validation manifest:
	lula dev print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

To print resources from lula validation manifest to output file:
	lula dev print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --output-file /path/to/output.json

To print the lula validation that generated a given observation:
	lula dev print --validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

```

### Options

```
  -a, --assessment string         the path to an assessment-results file
  -c, --component string          the path to a validation manifest file
  -h, --help                      help for print
  -u, --observation-uuid string   the observation uuid
  -o, --output-file string        the path to write the resources json
  -r, --resources                 true if the user is printing resources
  -v, --validation                true if the user is printing validation
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula dev](./lula_dev.md)	 - Collection of dev commands to make dev life easier

