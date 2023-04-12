package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {
	t.Run("GetOperatingSystem", testGetOperatingSystem)
	t.Run("Analyse", testAnalyse)
}

func testGetOperatingSystem(t *testing.T) {
	operatingSystem := GetOperatingSystem()
	assert.NotEqual(t, UnknownOperatingSystem, operatingSystem, "The operating system should not be unknown.")
}

func testAnalyse(t *testing.T) {
	analysis := Analyse()

	require.NotNil(t, analysis, "Analysis should not be nil.")
	assert.NotEmpty(t, analysis.GetApplicationName(), "Application name should not be empty.")
	assert.NotEqual(t, UnknownOperatingSystem, analysis.GetOperatingSystem(), "The operating system should not be unknown.")
	assert.NotEmpty(t, analysis.GetHomeDirectory(), "Home directory should not be empty.")
	assert.NotEmpty(t, analysis.GetUsername(), "Username should not be empty.")
	assert.NotEmpty(t, analysis.GetConfigFile(), "Config file should not be empty.")
}
