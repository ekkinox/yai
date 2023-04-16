package ui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUiInput(t *testing.T) {
	t.Run("NewUIInput", testNewUIInput)
	t.Run("GetRunMode", testGetRunMode)
	t.Run("GetPromptMode", testGetPromptMode)
	t.Run("GetArgs", testGetArgs)
}

func testNewUIInput(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-e"}
	uiInput, err := NewUIInput()
	assert.NoError(t, err, "NewUIInput should not return an error.")
	assert.NotNil(t, uiInput, "UiInput should not be nil.")
}

func testGetRunMode(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "arg1", "arg2"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, CliMode, uiInput.GetRunMode(), "RunMode should be CliMode.")
}

func testGetPromptMode(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-e"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, ExecPromptMode, uiInput.GetPromptMode(), "PromptMode should be ExecPromptMode.")
}

func testGetArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "arg1", "arg2"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, "arg1 arg2", uiInput.GetArgs(), "Args should be 'arg1 arg2'.")
}
