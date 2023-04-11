package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngineExecOutputGetCommand(t *testing.T) {
	eo := EngineExecOutput{Command: "testCommand"}
	result := eo.GetCommand()

	assert.Equal(t, "testCommand", result)
}

func TestEngineExecOutputGetExplanation(t *testing.T) {
	eo := EngineExecOutput{Explanation: "testExplanation"}
	result := eo.GetExplanation()

	assert.Equal(t, "testExplanation", result)
}

func TestEngineExecOutputIsExecutable(t *testing.T) {
	eo := EngineExecOutput{Executable: true}
	result := eo.IsExecutable()

	assert.True(t, result)
}

func TestEngineChatStreamOutputGetContent(t *testing.T) {
	co := EngineChatStreamOutput{content: "testContent"}
	result := co.GetContent()

	assert.Equal(t, "testContent", result)
}

func TestEngineChatStreamOutputIsLast(t *testing.T) {
	co := EngineChatStreamOutput{last: true}
	result := co.IsLast()

	assert.True(t, result)
}

func TestEngineChatStreamOutputIsInterrupt(t *testing.T) {
	co := EngineChatStreamOutput{interrupt: true}
	result := co.IsInterrupt()

	assert.True(t, result)
}

func TestEngineChatStreamOutputIsExecutable(t *testing.T) {
	co := EngineChatStreamOutput{executable: true}
	result := co.IsExecutable()

	assert.True(t, result)
}
