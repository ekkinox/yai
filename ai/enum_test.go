package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngineModeString(t *testing.T) {
	tests := []struct {
		name     string
		mode     EngineMode
		expected string
	}{
		{
			name:     "ExecEngineMode",
			mode:     ExecEngineMode,
			expected: "exec",
		},
		{
			name:     "ChatEngineMode",
			mode:     ChatEngineMode,
			expected: "chat",
		},
		{
			name:     "UnknownEngineMode",
			mode:     EngineMode(42),
			expected: "chat",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.mode.String()
			assert.Equal(t, test.expected, result)
		})
	}
}
