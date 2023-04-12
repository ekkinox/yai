package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestUIPrompt(t *testing.T) {
	t.Run("Prompt", testPrompt)
	t.Run("PromptStyle", testPromptStyle)
	t.Run("PromptIcon", testPromptIcon)
	t.Run("PromptPlaceholder", testPromptPlaceholder)
}

func testPrompt(t *testing.T) {
	testCases := []struct {
		name         string
		mode         PromptMode
		initialValue string
	}{
		{"Exec", ExecPromptMode, ""},
		{"Config", ConfigPromptMode, ""},
		{"Chat", ChatPromptMode, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPrompt(tc.mode)
			assert.Equal(t, tc.mode, p.GetMode(), "The prompt mode should match the expected value.")
			assert.Equal(t, tc.initialValue, p.GetValue(), "The prompt value should match the expected value.")
		})
	}
}

func testPromptStyle(t *testing.T) {
	testCases := []struct {
		name      string
		mode      PromptMode
		styleFunc func(PromptMode) lipgloss.Style
	}{
		{"Exec", ExecPromptMode, getPromptStyle},
		{"Config", ConfigPromptMode, getPromptStyle},
		{"Chat", ChatPromptMode, getPromptStyle},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			style := tc.styleFunc(tc.mode)
			assert.NotNil(t, style, "The prompt style should not be nil.")
		})
	}
}

func testPromptIcon(t *testing.T) {
	testCases := []struct {
		name     string
		mode     PromptMode
		iconFunc func(PromptMode) string
	}{
		{"Exec", ExecPromptMode, getPromptIcon},
		{"Config", ConfigPromptMode, getPromptIcon},
		{"Chat", ChatPromptMode, getPromptIcon},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			icon := tc.iconFunc(tc.mode)
			assert.NotEmpty(t, icon, "The prompt icon should not be empty.")
		})
	}
}

func testPromptPlaceholder(t *testing.T) {
	testCases := []struct {
		name            string
		mode            PromptMode
		placeholderFunc func(PromptMode) string
	}{
		{"Exec", ExecPromptMode, getPromptPlaceholder},
		{"Config", ConfigPromptMode, getPromptPlaceholder},
		{"Chat", ChatPromptMode, getPromptPlaceholder},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			placeholder := tc.placeholderFunc(tc.mode)
			assert.NotEmpty(t, placeholder, "The prompt placeholder should not be empty.")
		})
	}
}
