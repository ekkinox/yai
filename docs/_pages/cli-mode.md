---
title: "CLI mode"
classes: wide
permalink: /cli-mode/
---

```shell
yo list all processes listening on port 8080
```

This will perform a single run, using your [preferred](#configuration) prompt mode.

```shell
yo -e show the disk usage of my docker resources
```

This will perform a single command line generation (enforcing `🚀 exec` prompt mode usage with `-e`).

```shell
yo -c generate me a go application example using fiber
```

This will reply to a single question (enforcing `💬 chat` prompt mode usage with `-c`).

```shell
cat some_script.go | yo -c generate unit tests
```

You can also `pipe` input that will be taken into account in your request.
