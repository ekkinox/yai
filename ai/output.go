package ai

type EngineExecOutput struct {
	Command     string `json:"cmd"`
	Explanation string `json:"exp"`
	Executable  bool   `json:"exec"`
}

func (eo EngineExecOutput) GetCommand() string {
	return eo.Command
}

func (eo EngineExecOutput) GetExplanation() string {
	return eo.Explanation
}

func (eo EngineExecOutput) IsExecutable() bool {
	return eo.Executable
}

type EngineChatOutput struct {
	content    string
	last       bool
	interrupt  bool
	executable bool
}

func (co EngineChatOutput) GetContent() string {
	return co.content
}

func (co EngineChatOutput) IsLast() bool {
	return co.last
}

func (co EngineChatOutput) IsInterrupt() bool {
	return co.interrupt
}

func (co EngineChatOutput) IsExecutable() bool {
	return co.executable
}
