package engine

type EngineMode int

const (
	ChatEngineMode EngineMode = iota
	RunEngineMode
)

func (m EngineMode) String() string {
	if m == ChatEngineMode {
		return "chat"
	} else {
		return "run"
	}
}

func EngineModeFromString(mode string) EngineMode {
	if mode == "run" {
		return RunEngineMode
	} else {
		return ChatEngineMode
	}
}
