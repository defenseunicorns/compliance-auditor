---
title: lula generate profile
description: Lula CLI command reference for <code>lula generate profile</code>.
type: docs
---
## lula generate profile

Generate a profile OSCAL template

### Synopsis

Generation of a Profile OSCAL artifact with controls included or excluded from a source catalog/profile.

```
lula generate profile [flags]
```

### Examples

```

To generate a profile with included controls:
	lula generate profile -s <catalog/profile source> -i ac-1,ac-2,ac-3

To specify the name and filetype of the generated artifact:
	lula generate profile -s <catalog/profile source> -i ac-1,ac-2,ac-3 -o my_profile.yaml

To generate a profile that includes all controls except a list specified controls:
	lula generate profile -s <catalog/profile source> -e ac-1,ac-2,ac-3

```

### Options

```
  -a, --all                  Include all controls from the source catalog/profile
  -e, --exclude strings      comma delimited list of controls to exclude from the source catalog/profile
  -h, --help                 help for profile
  -i, --include strings      comma delimited list of controls to include from the source catalog/profile
  -o, --output-file string   the path to the output file. If not specified, the output file will be directed to stdout
  -s, --source string        the path to the source catalog/profile
```

### Options inherited from parent commands

```
  -f, --input-file string   Path to a manifest file
  -l, --log-level string    Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula generate](./lula_generate.md)	 - Generate a specified compliance artifact template

