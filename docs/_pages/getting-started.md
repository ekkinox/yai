---
title: "Getting started"
classes: wide
permalink: /getting-started/
---

## What is `Yai` ?

`Yai` (your AI) is an assistant for your terminal, unleashing the power of artificial intelligence to streamline your command line experience.

It is already aware of your:
- operating system & distribution
- username, shell & home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your experience.

## Installation

To install `Yai`, simply run:

```shell
curl -sS https://raw.githubusercontent.com/ekkinox/yai/main/install.sh | bash
```

- this will detect the proper binary to install for your machine
- and upgrade to the latest stable version if already installed

You can also install it from the [available releases](https://github.com/ekkinox/yai/releases) from the GitHub repository.

To uninstall `Yai`, simply run the provided uninstallation script:

```shell
curl -sS https://raw.githubusercontent.com/ekkinox/yai/main/uninstall.sh | bash
```

## Configuration

At first run, `Yai` will ask you to provide an [OpenAI API key](https://platform.openai.com/account/api-keys) (required to interact with **ChatGPT AI**).

It will then generate your configuration in the file `~/.config/yai.json`, with the following structure:

```json
{
  "openai_key": "sk-xxxxxxxxx",       // OpenAI API key (mandatory)
  "openai_model": "gpt-3.5-turbo",    // OpenAI API model (default gpt-3.5-turbo)
  "openai_proxy": "",                 // OpenAI API proxy (default disabled)
  "openai_temperature": 0.2,          // OpenAI API temperature (defaut 0.2)
  "openai_max_tokens": 1000,          // OpenAI API max tokens (default 1000)
  "user_default_prompt_mode": "exec", // user prefered prompt mode: "exec" (default) or "chat"
  "user_preferences": ""              // user preferences, expressed in natural language (default none)
}
```

## Fine tuning

You can fine tune `Yai` by editing the settings in the `~/.config/yai.json` configuration file.

Note that in `REPL` mode, you can press anytime `ctrl+s` to edit settings:
- it will open your editor on the configuration file
- and will hot reload settings changes when you're done.

### Model 

You can use the `openai_model` to configure the AI model you want to use. By default, the model `gpt-3.5-turbo` is used.

```json
{
  "openai_model": "gpt-4"
}
```

You can find the list of [supported models here](https://platform.openai.com/docs/models/model-endpoint-compatibility) (must be compatible with OpenAI API `v1/chat/completions`).

### Preferences

You can use the `user_preferences` to express any preferences in your natural language:

```json
{
  "user_preferences": "I want you to talk like an humorist, and I want you to always add the -y flag when I use dnf"
}
```

`Yai` will take them into account.