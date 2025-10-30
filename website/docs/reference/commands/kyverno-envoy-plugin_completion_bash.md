---
title: "kyverno-envoy-plugin completion bash"
slug: "kyverno-envoy-plugin_completion_bash"
description: "CLI reference for kyverno-envoy-plugin completion bash"
---

## kyverno-envoy-plugin completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(kyverno-envoy-plugin completion bash)

To load completions for every new session, execute once:

#### Linux:

	kyverno-envoy-plugin completion bash > /etc/bash_completion.d/kyverno-envoy-plugin

#### macOS:

	kyverno-envoy-plugin completion bash > $(brew --prefix)/etc/bash_completion.d/kyverno-envoy-plugin

You will need to start a new shell for this setup to take effect.


```
kyverno-envoy-plugin completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [kyverno-envoy-plugin completion](kyverno-envoy-plugin_completion.md)	 - Generate the autocompletion script for the specified shell

