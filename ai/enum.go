package ai

type EngineMode int

const (
	ExecEngineMode EngineMode = iota
	ChatEngineMode
)

func (m EngineMode) String() string {
	if m == ExecEngineMode {
		return "exec"
	} else {
		return "chat"
	}
}

func FromString(s string) EngineMode {
	if s == "exec" {
		return ExecEngineMode
	} else {
		return ChatEngineMode
	}
}
