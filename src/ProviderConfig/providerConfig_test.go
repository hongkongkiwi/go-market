package ProviderConfig

import (
	"os"
	"testing"

	"github.com/hongkongkiwi/chaostheory/src/Helpers"
)

func TestProviderConfig(t *testing.T) {
	// Test OpenDB and CloseDB functions
	t.Run("TestOpenCloseDB", func(t *testing.T) {
		tmpDBFileName, tempErr := Helpers.CreateTempFile("TestProviderFunctions")
		if tempErr != nil {
			t.Errorf("Error creating temporary file: %v", tempErr)
			return
		}
		err := OpenDB(tmpDBFileName.Name())
		if err != nil {
			t.Errorf("Error opening database: %v", err)
		}
		defer func() {
			CloseDB()
			os.Remove(tmpDBFileName.Name())
		}()
	})

	// Test SetProvider, GetProvider, and GetAllProviders functions
	t.Run("TestProviderFunctions", func(t *testing.T) {
		tmpDBFileName, tempErr := Helpers.CreateTempFile("TestProviderFunctions")
		if tempErr != nil {
			t.Errorf("Error creating temporary file: %v", tempErr)
			return
		}
		err := OpenDB(tmpDBFileName.Name())
		if err != nil {
			t.Errorf("Error opening database: %v", err)
		}
		defer func() {
			CloseDB()
			os.Remove(tmpDBFileName.Name())
		}()

		provider := &Provider{
			Name:  "TestProvider",
			Pairs: map[string]bool{"BTC/USD": true, "ETH/USD": false},
		}

		// Test SetProvider
		err = SetProvider(provider)
		if err != nil {
			t.Errorf("Error setting provider: %v", err)
		}

		// Test GetProvider
		retProvider, err := GetProvider(provider.Name)
		if err != nil {
			t.Errorf("Error getting provider: %v", err)
		}
		if retProvider == nil || retProvider.Name != provider.Name || retProvider.Pairs["BTC/USD"] != provider.Pairs["BTC/USD"] || retProvider.Pairs["ETH/USD"] != provider.Pairs["ETH/USD"] {
			t.Errorf("Retrieved provider does not match expected provider")
		}

		// Test GetAllProviders
		providers, err := GetProviders()
		if err != nil {
			t.Errorf("Error getting all providers: %v", err)
		}
		if len(providers) == 0 {
			t.Errorf("No providers retrieved")
		}
	})

}
