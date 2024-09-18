---
title: lula validate
description: Lula CLI command reference for <code>lula validate</code>.
type: docs
---
## lula validate

validate an OSCAL component definition

### Synopsis

Lula Validation of an OSCAL component definition

```
lula validate [flags]
```

### Examples

```

To validate on a cluster:
	lula validate -f ./oscal-component.yaml
To indicate a specific Assessment Results file to create or append to:
	lula validate -f ./oscal-component.yaml -o assessment-results.yaml
To target a specific control-implementation source / standard/ framework
	lula validate -f ./oscal-component.yaml -t critical
To run validations and automatically confirm execution
	lula dev validate -f ./oscal-component.yaml --confirm-execution
To run validations non-interactively (no execution)
	lula dev validate -f ./oscal-component.yaml --non-interactive

```

### Options

```
      --confirm-execution    confirm execution scripts run as part of the validation
  -h, --help                 help for validate
  -f, --input-file string    the path to the target OSCAL component definition
      --non-interactive      run the command non-interactively
  -o, --output-file string   the path to write assessment results. Creates a new file or appends to existing files
  -t, --target string        the specific control implementations or framework to validate against
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula](./lula.md)	 - Risk Management as Code

