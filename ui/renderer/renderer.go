package renderer

import (
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const error_color = "#cc3333"
const warning_color = "#ffcc00"
const success_color = "#46b946"

type Renderer struct {
	contentRenderer *glamour.TermRenderer
	successRenderer lipgloss.Style
	warningRenderer lipgloss.Style
	errorRenderer   lipgloss.Style
}

func NewRenderer(options ...glamour.TermRendererOption) *Renderer {

	contentRenderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		log.Printf("error: %v", err)
	}

	successRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(success_color))
	warningRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(warning_color))
	errorRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(error_color))

	return &Renderer{
		contentRenderer: contentRenderer,
		successRenderer: successRenderer,
		warningRenderer: warningRenderer,
		errorRenderer:   errorRenderer,
	}
}

func (r *Renderer) RenderContent(in string) string {
	out, _ := r.contentRenderer.Render(in)

	return out
}

func (r *Renderer) RenderSuccess(in string) string {
	return r.successRenderer.Render(in)
}

func (r *Renderer) RenderWarning(in string) string {
	return r.warningRenderer.Render(in)
}

func (r *Renderer) RenderError(in string) string {
	return r.errorRenderer.Render(in)
}
