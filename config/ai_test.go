package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAiConfig(t *testing.T) {
	t.Run("GetKey", testGetKey)
	t.Run("GetProxy", testGetProxy)
	t.Run("GetTemperature", testGetTemperature)
}

func testGetKey(t *testing.T) {
	expectedKey := "test_key"
	aiConfig := AiConfig{key: expectedKey}

	actualKey := aiConfig.GetKey()

	assert.Equal(t, expectedKey, actualKey, "The two keys should be the same.")
}

func testGetProxy(t *testing.T) {
	expectedProxy := "test_proxy"
	aiConfig := AiConfig{proxy: expectedProxy}

	actualProxy := aiConfig.GetProxy()

	assert.Equal(t, expectedProxy, actualProxy, "The two proxies should be the same.")
}

func testGetTemperature(t *testing.T) {
	expectedTemperature := 0.7
	aiConfig := AiConfig{temperature: expectedTemperature}

	actualTemperature := aiConfig.GetTemperature()

	assert.Equal(t, expectedTemperature, actualTemperature, "The two temperatures should be the same.")
}
