# Yo

> Your AI powered CLI assistant.

## Installation

```shell
go get && go build -o /usr/local/bin/yo && sudo chmod +x /usr/local/bin/yo
```

## Usage

### UI modes

Yo provides 2 UI modes:
- TUI: terminal user interface, made to offer interactive prompts like a discussion
- CLI: command line interface, made to run a single execution

#### TUI mode

```shell
yo
```

This will open a REPL loop, where you can use the following shortcuts:

| Keys     | Description                                   |
|----------|-----------------------------------------------|
| `↑` `↓`  | Navigate in history                           |
| `tab`    | Switch between `💬 chat` or `🚀 run` modes     |
| `ctrl+r` | Clear terminal and reset discussion history   |
| `ctrl+l` | Clear terminal but keep discussion history    |
| `ctrl+s` | Edit settings                                 |
| `ctrl+c` | Exit or interrupt current command / completion |


#### CLI mode

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
  "user_default_mode": "chat",  // prefered run mode: [chat] or [run]
  "user_context": ""            // to express some preferences in natural language
}
```