package components

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func NewPrompt() textinput.Model {

	prompt := textinput.New()
	prompt.Placeholder = "Enter something"
	prompt.Focus()

	return prompt
}
