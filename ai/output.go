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

type EngineChatStreamOutput struct {
	content    string
	last       bool
	interrupt  bool
	executable bool
}

func (co EngineChatStreamOutput) GetContent() string {
	return co.content
}

func (co EngineChatStreamOutput) IsLast() bool {
	return co.last
}

func (co EngineChatStreamOutput) IsInterrupt() bool {
	return co.interrupt
}

func (co EngineChatStreamOutput) IsExecutable() bool {
	return co.executable
}
