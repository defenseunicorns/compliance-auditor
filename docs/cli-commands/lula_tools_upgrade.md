---
title: lula tools upgrade
description: Lula CLI command reference for <code>lula tools upgrade</code>.
type: docs
---
## lula tools upgrade

Upgrade OSCAL document to a new version if possible.

### Synopsis

Validate an OSCAL document against the OSCAL schema version provided. If the document is valid, upgrade it to the provided OSCAL version. Otherwise, return or write as ValidationError. Yaml formatting handled by gopkg/yaml.v3 and while objects will maintain deep equality, visual representation may be different than the input file.

```
lula tools upgrade [flags]
```

### Examples

```

To Upgrade an existing OSCAL file:
	lula tools upgrade -f <path to oscal> -v <version>

```

### Options

```
  -h, --help                       help for upgrade
  -f, --input-file string          the path to a oscal json schema file
  -o, --output-file string         the path to write the linted oscal json schema file (default is the input file)
  -r, --validation-result string   the path to write the validation result file
  -v, --version string             the version of the oscal schema to validate against (default is the latest supported version) (default "1.1.2")
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](./lula_tools.md)	 - Collection of additional commands to make OSCAL easier

