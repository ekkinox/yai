package prompt

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
