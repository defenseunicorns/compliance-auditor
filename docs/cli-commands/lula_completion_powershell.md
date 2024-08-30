---
title: lula completion powershell
description: Lula CLI command reference for <code>lula completion powershell</code>.
type: docs
---
## lula completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	lula completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
lula completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula completion](/cli/cli-commands/lula_completion/)	 - Generate the autocompletion script for the specified shell

