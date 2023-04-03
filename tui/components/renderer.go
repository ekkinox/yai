package components

import (
	"github.com/charmbracelet/glamour"
	"log"
)

type Renderer struct {
	renderer *glamour.TermRenderer
}

func NewRenderer(options ...glamour.TermRendererOption) *Renderer {

	renderer, err := glamour.NewTermRenderer(options...)

	if err != nil {
		log.Printf("error: %v", err)
	}

	return &Renderer{renderer}
}

func (r *Renderer) Render(in string) string {

	// allow silent failures for text streams rendering
	out, _ := r.renderer.Render(in)

	return out
}
