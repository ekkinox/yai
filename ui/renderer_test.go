package ui

import (
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/stretchr/testify/assert"
)

func TestUIRenderer(t *testing.T) {
	t.Run("Renderer", testRenderer)
	t.Run("RenderContent", testRenderContent)
	t.Run("RenderSuccess", testRenderSuccess)
	t.Run("RenderWarning", testRenderWarning)
	t.Run("RenderError", testRenderError)
	t.Run("RenderHelp", testRenderHelp)
	t.Run("RenderConfigMessage", testRenderConfigMessage)
	t.Run("RenderHelpMessage", testRenderHelpMessage)
}

func testRenderer(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	assert.NotNil(t, r, "Renderer should not be nil.")
}

func testRenderContent(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	input := "Hello, World!"
	output := r.RenderContent(input)
	assert.NotEmpty(t, output, "Rendered content should not be empty.")
}

func testRenderSuccess(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	input := "Success message"
	output := r.RenderSuccess(input)
	assert.NotEmpty(t, output, "Rendered success message should not be empty.")
}

func testRenderWarning(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	input := "Warning message"
	output := r.RenderWarning(input)
	assert.NotEmpty(t, output, "Rendered warning message should not be empty.")
}

func testRenderError(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	input := "Error message"
	output := r.RenderError(input)
	assert.NotEmpty(t, output, "Rendered error message should not be empty.")
}

func testRenderHelp(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	input := "Help message"
	output := r.RenderHelp(input)
	assert.NotEmpty(t, output, "Rendered help message should not be empty.")
}

func testRenderConfigMessage(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	output := r.RenderConfigMessage()
	assert.NotEmpty(t, output, "Rendered config message should not be empty.")
}

func testRenderHelpMessage(t *testing.T) {
	r := NewRenderer(glamour.WithAutoStyle())
	output := r.RenderHelpMessage()
	assert.NotEmpty(t, output, "Rendered help message should not be empty.")
}
