# 🚀 Yo 💬 - AI powered terminal assistant

[![build](https://github.com/ekkinox/yo/actions/workflows/build.yml/badge.svg)](https://github.com/ekkinox/yo/actions/workflows/build.yml)
[![release](https://github.com/ekkinox/yo/actions/workflows/release.yml/badge.svg)](https://github.com/ekkinox/yo/actions/workflows/release.yml)
[![doc](https://github.com/ekkinox/yo/actions/workflows/doc.yml/badge.svg)](https://github.com/ekkinox/yo/actions/workflows/doc.yml)

> Unleash the power of artificial intelligence to streamline your command line experience.

![Intro](docs/_assets/intro.gif)

## What is Yo ?

`Yo` is an assistant for your terminal, using [OpenAI ChatGPT](https://chat.openai.com/) to build and run commands for you. You just need to describe them in your everyday language, it will take care or the rest. 

You have any questions on random topics in mind? You can also ask `Yo`, and get the power of AI without leaving `/home`.

It is already aware of your:
- operating system & distribution
- username, shell & home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your experience.

## Documentation

A complete documentation is available at [https://ekkinox.github.io/yo/](https://ekkinox.github.io/yo/).

## Quick start

To install `Yo`, simply run:

```shell
curl -sS https://raw.githubusercontent.com/ekkinox/yo/main/install.sh | bash
```

You can also install it from the [available releases](https://github.com/ekkinox/yo/releases).

At first run, it will ask you for an [OpenAI API key](https://platform.openai.com/account/api-keys), and use it to create the configuration file in `~/.config/yo.json`.

See [documentation](https://ekkinox.github.io/yo/getting-started/#configuration) for more information.



