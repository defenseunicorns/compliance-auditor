---
title: lula tools uuidgen
description: Lula CLI command reference for <code>lula tools uuidgen</code>.
type: docs
---
## lula tools uuidgen

Generate a UUID

### Synopsis

Generate a UUID at random or deterministically with a provided string

```
lula tools uuidgen [flags]
```

### Examples

```

To create a new random UUID:
	lula tools uuidgen

To create a deterministic UUID given some source:
	lula tools uuidgen <source>

```

### Options

```
  -h, --help   help for uuidgen
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula tools](/cli/cli-commands/lula_tools/)	 - Collection of additional commands to make OSCAL easier

