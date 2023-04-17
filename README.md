# 🚀 Yo 💬

> AI powered terminal assistant.

![Intro](docs/examples/intro.gif)

## Table of contents

<!-- TOC -->
* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
  * [CLI mode](#cli-mode)
  * [REPL mode](#repl-mode)
* [Configuration](#configuration)
* [Examples](#examples)
* [Changelog](#changelog)
<!-- TOC -->

## Description

`Yo` is an AI powered assistant for your terminal, unleashing the power of AI to streamline your command line experience.

You can also engage in conversation with it to get help on any topics.

It is already aware of your:
- operating system
- distribution
- username
- shell
- home directory
- preferred editor

You can also give any supplementary preferences to fine tune your experience (see [configuration](#configuration) section).
## Installation

```shell
curl -sS https://raw.githubusercontent.com/ekkinox/yo/main/install.sh | bash
```

You can also install it from the [available releases](https://github.com/ekkinox/yo/releases).

## Usage

`Yo` provides 2 run modes:
- `CLI` mode: command line interface, made to perform a single run
- `REPL` mode: terminal user interface, made to offer interactive prompts in a loop

### CLI mode

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

### REPL mode

```shell
yo
```

This will open a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface, with 2 types of prompts

- `🚀 exec`: will generate a command line to execute for what you're asking
- `💬 chat`: will engage in a discussion to help you the best way possible

You also can use the following `keyboard shortcuts`:

| Keys     | Description                                         |
|----------|-----------------------------------------------------|
| `↑` `↓`  | Navigate in history                                 |
| `tab`    | Switch between `🚀 exec` and `💬 chat` prompt modes |
| `ctrl+h` | Show help                                           |
| `ctrl+s` | Edit settings                                       |
| `ctrl+r` | Clear terminal and reset discussion history         |
| `ctrl+l` | Clear terminal but keep discussion history          |
| `ctrl+c` | Exit or interrupt command execution                 |


## Configuration

At the first execution, `Yo` will ask you to provide an [OpenAI API key](https://platform.openai.com/account/api-keys).

It will then generate your configuration in the file `~/.config/yo.json`, with the following structure:

```JS
{
  "openai_key": "sk-xxxxxxxxx",       // OpenAI API key
  "openai_proxy": "",                 // OpenAI API proxy (if needed)
  "openai_temperature": 0.2,          // OpenAI API temperature
  "user_default_prompt_mode": "exec", // user prefered prompt mode: "exec" (default) or "chat"
  "user_preferences": ""              // user preferences, expressed in natural language
}
```

In `REPL` mode, you can press `ctrl+s` to edit settings, they will be hot reloaded once you close your editor.

## Examples

Check the [examples](/examples) folder for some use case ideas.

## Changelog

See [CHANGELOG](CHANGELOG.md).