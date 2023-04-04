# Yo

> Your AI powered CLI assistant.

## Installation

```shell
sudo go build -o /usr/local/bin/yo && sudo chmod +x /usr/local/bin/yo
```

At the first execution, your assistant will ask you to provide an [OpenAI API key](https://platform.openai.com/account/api-keys).

## Usage

### Run modes

Yo provides 2 run modes:
- TUI: terminal user interface, made to offer interactive prompts like a discussion
- CLI: command line interface, made to run a single execution

Launch with TUI run mode:
```shell
yo
```

Launch with CLI run mode:
```shell
yo list all my files in my home directory
```

### Keyboard shortcuts

| Keys     | Description                                     |
|----------|-------------------------------------------------|
| `↑` `↓`  | Navigate history                                |
| `tab`    | Switch between `chat` or `run` modes            |
| `ctrl+l` | Clear terminal but keep discussion history      |
| `ctrl+s` | Edit settings                                   |
| `ctrl+c` | Exit or interrupt current command or completion |


