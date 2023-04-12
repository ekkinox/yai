# 🚀 / 💬 Yo 

> AI powered terminal assistant.

![Demo](doc/demo.gif)

## Table of contents

<!-- TOC -->
* [Description](#description)
* [Installation](#installation)
* [Configuration](#configuration)
* [Usage](#usage)
  * [REPL mode](#repl-mode)
  * [CLI mode](#cli-mode)
<!-- TOC -->

## Description

`Yo` is an AI powered assistant for your terminal.

You can converse with it using your natural language, and it will provide you with ready to use command lines, or engage in a discussion to assist you.

`Yo` is already aware of your:
- operating system
- distribution
- username
- shell
- home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your user experience (see [configuration](#configuration) section).

## Installation

```shell
go get && sudo go build -o /usr/local/bin/yo && sudo chmod +x /usr/local/bin/yo
```

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

## Usage

`Yo` provides 2 run modes:
- **REPL mode**: terminal user interface, made to offer interactive prompts in a loop
- **CLI mode**: command line interface, made to perform a single run

### REPL mode

```shell
yo
```

This will open a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface, with 2 types of prompts

- `🚀 exec`: will generate a command line to execute for what you're asking
- `💬 chat`: will engage in a discussion to help you the best way possible

You also can use the following `keyboard shortcuts`:

| Keys     | Description                                           |
|----------|-------------------------------------------------------|
| `↑` `↓`  | Navigate in history                                   |
| `tab`    | Switch between `🚀 exec` and `💬 chat` prompt modes   |
| `ctrl+s` | Open editor on configuration file                     |
| `ctrl+r` | Clear terminal and reset discussion history           |
| `ctrl+l` | Clear terminal but keep discussion history            |
| `ctrl+c` | Exit or interrupt command execution / chat completion |


### CLI mode

```shell
yo list all javascript files
```

This will perform a single run, using your [preferred prompt mode](#configuration).

```shell
yo -e list all javascript files
```

This will perform a single command line generation, enforcing `🚀 exec` prompt mode usage with `-e`.

```shell
yo -c how can I find javascript files
```

This will reply to a single question, enforcing `💬 chat` prompt mode usage with `-c`.