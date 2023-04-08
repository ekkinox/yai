# Yo

> Your AI powered terminal assistant.

![Demo](doc/demo.gif)

## Table of contents

<!-- TOC -->
* [Description](#description)
* [Installation](#installation)
* [Configuration](#configuration)
* [Usage](#usage)
  * [TUI mode](#tui-mode)
  * [CLI mode](#cli-mode)
<!-- TOC -->

## Description

`Yo` is an AI powered assistant for your terminal.

You can converse in your natural language, and it will provide you with either pre-built command lines or engage in discussion to assist you.

To be the most accurate possible, it is already aware of your:
- operating system
- distribution
- username
- shell
- home directory
- preferred editor

You can also give any preferences to fine tune your user experience (see [configuration](#configuration) section).

## Installation

```shell
go get && sudo go build -o /usr/local/bin/yo && sudo chmod +x /usr/local/bin/yo
```

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

## Usage

`Yo` provides 2 UI modes:
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
| `ctrl+s` | Edit settings                                  |
| `ctrl+r` | Clear terminal and reset discussion history    |
| `ctrl+l` | Clear terminal but keep discussion history     |
| `ctrl+c` | Exit or interrupt current command / completion |


### CLI mode

**TODO**: work in progress

```shell
yo list all my files in my home directory
```

This will perform a single execution, according to your input.

