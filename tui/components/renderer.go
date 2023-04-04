package components

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"log"
)

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

	successRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color("#46b946"))
	warningRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00"))
	errorRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color("#cc3333"))

	return &Renderer{
		contentRenderer: contentRenderer,
		successRenderer: successRenderer,
		warningRenderer: warningRenderer,
		errorRenderer:   errorRenderer,
	}
}

func (r *Renderer) RenderContent(in string) string {

	// allow silent failures for text streams rendering
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
