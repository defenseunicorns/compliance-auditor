---
title: lula generate system-security-plan
description: Lula CLI command reference for <code>lula generate system-security-plan</code>.
type: docs
---
## lula generate system-security-plan

Generate a system security plan OSCAL artifact

### Synopsis

Generation of a System Security Plan OSCAL artifact from a source profile along with an optional list of component definitions.

```
lula generate system-security-plan [flags]
```

### Examples

```

To generate a system security plan from profile and component definition:
	lula generate system-security-plan -p <path/to/profile> -c <path/to/component-definition>

To specify the name and filetype of the generated artifact:
	lula generate system-security-plan -p <path/to/profile> -c <path/to/component-definition> -o my_ssp.yaml

```

### Options

```
  -c, --components strings                      comma delimited list the paths to the component definitions to include for the SSP
  -h, --help                                    help for system-security-plan
  -o, --output-file system-security-plan.yaml   the path to the output file. If not specified, the output file will default to system-security-plan.yaml
  -p, --profile string                          the path to the imported profile
      --remarks strings                         Target for remarks population (default [statement])
```

### Options inherited from parent commands

```
  -f, --input-file string   Path to a manifest file
  -l, --log-level string    Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula generate](./lula_generate.md)	 - Generate a specified compliance artifact template

