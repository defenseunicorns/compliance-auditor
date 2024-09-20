---
title: lula generate component
description: Lula CLI command reference for <code>lula generate component</code>.
type: docs
---
## lula generate component

Generate a component definition OSCAL template

```
lula generate component [flags]
```

### Examples

```

To generate a new component-definition template:
lula generate component -c <catalog source url> -r control-a,control-b,control-c
- IE lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-2,au-5

To Generate and merge with an existing Component Definition:
lula generate component -c <catalog source url> -r control-a,control-b,control-c -o existing-component.yaml

To Generate a component definition with a specific "named" component:
lula generate component -c <catalog source url> -r control-a --component "Software X"

To Generate a component definition with remarks populated from specific control "parts":
lula generate component -c <catalog source url> -r control-a --remarks guidance,assessment-objective

```

### Options

```
  -c, --catalog-source string   Catalog source location (local or remote)
      --component string        Component Title
      --framework string        Control-Implementation collection that these controls belong to
  -h, --help                    help for component
  -p, --profile string          Profile source location (local or remote)
      --remarks strings         Target for remarks population (default = statement)
  -r, --requirements strings    List of requirements to capture
```

### Options inherited from parent commands

```
  -f, --input-file string    Path to a manifest file
  -l, --log-level string     Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
  -o, --output-file string   Path and Name to an output file
```

### SEE ALSO

* [lula generate](./lula_generate.md)	 - Generate a specified compliance artifact template

