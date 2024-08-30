---
title: lula completion zsh
description: Lula CLI command reference for <code>lula completion zsh</code>.
type: docs
---
## lula completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(lula completion zsh)

To load completions for every new session, execute once:

#### Linux:

	lula completion zsh > "${fpath[1]}/_lula"

#### macOS:

	lula completion zsh > $(brew --prefix)/share/zsh/site-functions/_lula

You will need to start a new shell for this setup to take effect.


```
lula completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula completion](/cli/cli-commands/lula_completion/)	 - Generate the autocompletion script for the specified shell

