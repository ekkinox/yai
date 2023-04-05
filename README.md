# Yo

> Your AI powered CLI assistant.

![Demo](doc/demo.gif)

## Table of contents

<!-- TOC -->
* [Installation](#installation)
* [Usage](#usage)
  * [TUI mode](#tui-mode)
  * [CLI mode](#cli-mode)
* [Configuration](#configuration)
<!-- TOC -->

## Installation

```shell
go get && sudo go build -o /usr/local/bin/yo && sudo chmod +x /usr/local/bin/yo
```

## Usage


Yo provides 2 UI modes:
- TUI: terminal user interface, made to offer interactive prompts like a discussion
- CLI: command line interface, made to perform a single execution

### TUI mode

```shell
yo
```

This will open a [REPL loop](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop), with 2 types of prompts

- `🚀 run`: will try to provide a command line for what you ask to do
- `💬 chat`: will chat with you to help you the best way possible

You also can use the following shortcuts:

| Keys     | Description                                    |
|----------|------------------------------------------------|
| `↑` `↓`  | Navigate in history                            |
| `tab`    | Switch between `🚀 run` and `💬 chat` prompts  |
| `ctrl+r` | Clear terminal and reset discussion history    |
| `ctrl+l` | Clear terminal but keep discussion history     |
| `ctrl+s` | Edit settings                                  |
| `ctrl+c` | Exit or interrupt current command / completion |


### CLI mode

**TODO**: work in progress

```shell
yo list all my files in my home directory
```

This will perform a single execution, according to your input.

## Configuration

At the first execution, your assistant will ask you to provide an [OpenAI API key](https://platform.openai.com/account/api-keys).

It will then generate your configuration in the file `~/.config/yo.json`, and will have the following structure:

```JS
{
  "openai_key": "sk-xxxxxxxxx", // your OpenAI API key
  "openai_temperature": 0.2,    // chatGPT temperature
  "user_default_mode": "run",   // prefered run mode: "run" (default) or "chat"
  "user_context": ""            // to express some preferences in natural language
}
```
