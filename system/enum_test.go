package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatingSystem(t *testing.T) {
	t.Run("String", testOperatingSystemString)
}

func testOperatingSystemString(t *testing.T) {
	testCases := []struct {
		name            string
		operatingSystem OperatingSystem
		expected        string
	}{
		{"Unknown", UnknownOperatingSystem, "unknown"},
		{"Linux", LinuxOperatingSystem, "linux"},
		{"macOS", MacOperatingSystem, "macOS"},
		{"Windows", WindowsOperatingSystem, "windows"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.operatingSystem.String(), "The string representation should match the expected value.")
		})
	}
}
