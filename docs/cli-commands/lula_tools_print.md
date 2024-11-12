---
title: lula tools print
description: Lula CLI command reference for <code>lula tools print</code>.
type: docs
---
## lula tools print

Print Resources or Lula Validation from an Assessment Observation

### Synopsis


Prints out data about an OSCAL Observation from the OSCAL Assessment Results model. 
Given "--resources", the command will print the JSON resources input that were provided to a Lula Validation, as identified by a given observation and assessment results file. 
Given "--validation", the command will print the Lula Validation that generated a given observation, as identified by a given observation, assessment results file, and component definition file.


```
lula tools print [flags]
```

### Examples

```

To print resources from lula validation manifest:
	lula tools print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

To print resources from lula validation manifest to output file:
	lula tools print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --output-file /path/to/output.json

To print the lula validation that generated a given observation:
	lula tools print --validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

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

* [lula tools](./lula_tools.md)	 - Collection of additional commands to make OSCAL easier

