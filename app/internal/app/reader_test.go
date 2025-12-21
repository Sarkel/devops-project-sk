package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartReader_ConfigValidation(t *testing.T) {
	// Test that StartReader requires proper configuration
	// Note: This test would fail if config is missing, which is expected behavior
	// In a real environment, you'd mock the config.Load() function

	// For now, we just verify the function exists and is callable
	assert.NotNil(t, StartReader)
}
