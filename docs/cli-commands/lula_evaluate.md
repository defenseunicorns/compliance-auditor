---
title: lula evaluate
description: Lula CLI command reference for <code>lula evaluate</code>.
type: docs
---
## lula evaluate

evaluate two results of a Security Assessment Results

### Synopsis

Lula evaluation of Security Assessment Results

```
lula evaluate [flags]
```

### Examples

```

To evaluate the latest results in two assessment results files:
	lula evaluate -f assessment-results-threshold.yaml -f assessment-results-new.yaml

To evaluate two results (threshold and latest) in a single OSCAL file:
	lula evaluate -f assessment-results.yaml

To target a specific framework for validation:
	lula evaluate -f assessment-results.yaml --target critical


```

### Options

```
  -h, --help                 help for evaluate
  -f, --input-file strings   Path to the file to be evaluated
  -s, --summary              Print a summary of the evaluation
  -t, --target string        the specific control implementations or framework to validate against
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula](./lula.md)	 - Risk Management as Code

