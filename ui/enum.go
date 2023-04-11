package ui

type PromptMode int

const (
	ExecPromptMode PromptMode = iota
	ConfigPromptMode
	ChatPromptMode
)

func (m PromptMode) String() string {
	switch m {
	case ExecPromptMode:
		return "exec"
	case ConfigPromptMode:
		return "config"
	default:
		return "chat"
	}
}

func FromString(s string) PromptMode {
	switch s {
	case "exec":
		return ExecPromptMode
	case "config":
		return ConfigPromptMode
	default:
		return ChatPromptMode
	}
}

type RunMode int

const (
	CliMode RunMode = iota
	ReplMode
)

func (m RunMode) String() string {
	if m == CliMode {
		return "cli"
	} else {
		return "repl"
	}
}
