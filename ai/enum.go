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
