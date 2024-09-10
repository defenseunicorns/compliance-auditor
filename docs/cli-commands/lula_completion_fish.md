---
title: lula completion fish
description: Lula CLI command reference for <code>lula completion fish</code>.
type: docs
---
## lula completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	lula completion fish | source

To load completions for every new session, execute once:

	lula completion fish > ~/.config/fish/completions/lula.fish

You will need to start a new shell for this setup to take effect.


```
lula completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula completion](./lula_completion.md)	 - Generate the autocompletion script for the specified shell

