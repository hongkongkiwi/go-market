package MarketSimulatorConfig

import (
	"testing"

	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/stretchr/testify/assert"
)

// TestGenerateDefaultConfig tests generating default configuration data.
func TestGenerateDefaultConfig(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestStartSimulation")
	if tempErr != nil {
		t.Errorf("Error creating temporary file: %v", tempErr)
		return
	}
	DbFile = tmpDBFileName.Name()

	// Generate default configuration
	defaultConfig := GenerateDefaultConfig()
	assert.NotNil(t, defaultConfig)
}
