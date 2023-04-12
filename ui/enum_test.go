package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUI(t *testing.T) {
	t.Run("PromptModeString", testPromptModeString)
	t.Run("GetPromptModeFromString", testGetPromptModeFromString)
	t.Run("RunModeString", testRunModeString)
}

func testPromptModeString(t *testing.T) {
	testCases := []struct {
		name       string
		promptMode PromptMode
		expected   string
	}{
		{"Exec", ExecPromptMode, "exec"},
		{"Config", ConfigPromptMode, "config"},
		{"Chat", ChatPromptMode, "chat"},
		{"Default", DefaultPromptMode, "default"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.promptMode.String(), "The string representation should match the expected value.")
		})
	}
}

func testGetPromptModeFromString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected PromptMode
	}{
		{"Exec", "exec", ExecPromptMode},
		{"Config", "config", ConfigPromptMode},
		{"Chat", "chat", ChatPromptMode},
		{"Default", "unknown", DefaultPromptMode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetPromptModeFromString(tc.input), "The prompt mode should match the expected value.")
		})
	}
}

func testRunModeString(t *testing.T) {
	testCases := []struct {
		name     string
		runMode  RunMode
		expected string
	}{
		{"CLI", CliMode, "cli"},
		{"REPL", ReplMode, "repl"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.runMode.String(), "The string representation should match the expected value.")
		})
	}
}
