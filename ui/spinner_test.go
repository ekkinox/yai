package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUISpinner(t *testing.T) {
	// Test cases
	t.Run("NewSpinner", testNewSpinner)
	t.Run("View", testSpinnerView)
}

func testNewSpinner(t *testing.T) {
	s := NewSpinner()
	assert.NotNil(t, s, "Spinner should not be nil.")
}

func testSpinnerView(t *testing.T) {
	s := NewSpinner()
	view := s.View()
	assert.NotEmpty(t, view, "Spinner view should not be empty.")
}
