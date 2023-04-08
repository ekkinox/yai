package prompt

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const exec_color = "#edc95e"
const exec_icon = "🚀 "
const exec_placeholder = "Ask me something..."
const config_color = "#ffffff"
const config_icon = "🔒 "
const config_placeholder = "Enter your OpenAI key..."
const chat_color = "#66b3ff"
const chat_icon = "💬 "
const chat_placeholder = "Ask me something..."

type Prompt struct {
	mode  PromptMode
	Input textinput.Model
}

func NewPrompt(mode PromptMode) *Prompt {

	input := textinput.New()
	input.Placeholder = getPromptPlaceholder(mode)
	input.TextStyle = getPromptStyle(mode)
	input.Prompt = getPromptIcon(mode)

	if mode == ConfigPromptMode {
		input.EchoMode = textinput.EchoPassword
	}

	input.Focus()

	return &Prompt{
		mode:  mode,
		Input: input,
	}

}

func (p *Prompt) GetMode() PromptMode {
	return p.mode
}

func (p *Prompt) ChangeMode(mode PromptMode) *Prompt {

	p.mode = mode

	p.Input.TextStyle = getPromptStyle(mode)
	p.Input.Prompt = getPromptIcon(mode)
	p.Input.Placeholder = getPromptPlaceholder(mode)

	return p
}

func (p *Prompt) String() string {
	style := getPromptStyle(p.mode)

	return fmt.Sprintf("%s%s\n", style.Render(getPromptIcon(p.mode)), style.Render(p.Input.Value()))
}

func getPromptStyle(mode PromptMode) lipgloss.Style {
	switch mode {
	case ExecPromptMode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(exec_color))
	case ConfigPromptMode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(config_color))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(chat_color))
	}
}

func getPromptIcon(mode PromptMode) string {
	switch mode {
	case ExecPromptMode:
		return exec_icon
	case ConfigPromptMode:
		return config_icon
	default:
		return chat_icon
	}
}

func getPromptPlaceholder(mode PromptMode) string {
	switch mode {
	case ExecPromptMode:
		return exec_placeholder
	case ConfigPromptMode:
		return config_placeholder
	default:
		return chat_placeholder
	}
}
