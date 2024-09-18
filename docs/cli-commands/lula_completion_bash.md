---
title: lula completion bash
description: Lula CLI command reference for <code>lula completion bash</code>.
type: docs
---
## lula completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(lula completion bash)

To load completions for every new session, execute once:

#### Linux:

	lula completion bash > /etc/bash_completion.d/lula

#### macOS:

	lula completion bash > $(brew --prefix)/etc/bash_completion.d/lula

You will need to start a new shell for this setup to take effect.


```
lula completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -l, --log-level string   Log level when running Lula. Valid options are: warn, info, debug, trace (default "info")
```

### SEE ALSO

* [lula completion](./lula_completion.md)	 - Generate the autocompletion script for the specified shell

