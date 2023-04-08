package ai

type EngineOutput struct {
	content    string
	last       bool
	interrupt  bool
	executable bool
}

func (d EngineOutput) GetContent() string {
	return d.content
}

func (d EngineOutput) IsLast() bool {
	return d.last
}

func (d EngineOutput) IsInterrupt() bool {
	return d.interrupt
}

func (d EngineOutput) IsExecutable() bool {
	return d.executable
}
