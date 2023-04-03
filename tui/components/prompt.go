package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekkinox/yo/engine"
)

func NewPrompt(mode engine.EngineMode) textinput.Model {

	prompt := textinput.New()
	prompt.Placeholder = "Ask me something..."
	prompt.TextStyle = GetPromptStyle(mode)
	prompt.Prompt = GetPromptIcon(mode)
	prompt.Focus()

	return prompt

}

func RenderPromptText(value string, mode engine.EngineMode) string {
	style := GetPromptStyle(mode)

	return fmt.Sprintf("\n%s%s\n", style.Render(GetPromptIcon(mode)), style.Render(value))
}

func UpdatePrompt(prompt textinput.Model, mode engine.EngineMode) textinput.Model {
	prompt.Prompt = GetPromptIcon(mode)
	prompt.TextStyle = GetPromptStyle(mode)

	return prompt
}

func GetPromptIcon(mode engine.EngineMode) string {
	if mode == engine.ChatEngineMode {
		return GetPromptStyle(mode).Render("[💬]> ")
	} else {
		return GetPromptStyle(mode).Render("[⚙️ ]> ")
	}
}

func GetPromptStyle(mode engine.EngineMode) lipgloss.Style {
	if mode == engine.ChatEngineMode {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#66b3ff"))
	} else {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#edc95e"))
	}
}
