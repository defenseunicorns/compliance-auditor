---
title: lula report
description: Lula CLI command reference for <code>lula report</code>.
type: docs
---
## lula report

Build a compliance report

```
lula report [flags]
```

### Examples

```

To create a new report:
lula report -f oscal-component-definition.yaml

To create a new report in json format:
lula report -f oscal-component-definition.yaml --file-format json

To create a new report in yaml format:
lula report -f oscal-component-definition.yaml --file-format yaml

```

### Options

```
      --file-format string   File format of the report (default "table")
  -h, --help                 help for report
  -f, --input-file string    Path to an OSCAL file
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula](./lula.md)	 - Risk Management as Code

