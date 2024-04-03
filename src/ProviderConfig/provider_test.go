package ProviderConfig

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/hongkongkiwi/chaostheory/src/Helpers"
)

func TestSerializeToBytes(t *testing.T) {
	provider := &Provider{
		Name: "TestProvider",
		Pairs: map[string]bool{
			"USD/EUR": true,
			"USD/GBP": false,
		},
	}

	expectedBytes, _ := json.Marshal(provider)

	bytes, err := provider.SerializeToBytes()
	if err != nil {
		t.Fatalf("Error serializing provider to bytes: %v", err)
	}

	if !reflect.DeepEqual(bytes, expectedBytes) {
		t.Errorf("Expected bytes %v; got %v", expectedBytes, bytes)
	}
}

func TestDeserializeProviderFromBytes(t *testing.T) {
	provider := &Provider{
		Name: "TestProvider",
		Pairs: map[string]bool{
			"USD/EUR": true,
			"USD/GBP": false,
		},
	}

	bytes, _ := json.Marshal(provider)

	deserializedProvider, err := DeserializeProviderFromBytes(bytes)
	if err != nil {
		t.Fatalf("Error deserializing provider from bytes: %v", err)
	}

	if !reflect.DeepEqual(deserializedProvider, provider) {
		t.Errorf("Expected deserialized provider %v; got %v", provider, deserializedProvider)
	}
}

func TestGetPairs(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestGetPairs")
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
		Name: "TestProvider",
		Pairs: map[string]bool{
			"USD/EUR": true,
			"USD/GBP": false,
		},
	}
	SetProvider(provider)

	pairs, err := provider.GetPairs(provider.Name)
	if err != nil {
		t.Fatalf("Error getting pairs: %v", err)
	}

	if !reflect.DeepEqual(pairs, provider.Pairs) {
		t.Errorf("Expected pairs %v; got %v", provider.Pairs, pairs)
	}
}

func TestIsPairEnabled(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestGetPairs")
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
		Name: "TestProvider",
		Pairs: map[string]bool{
			"USD/EUR": true,
			"USD/GBP": false,
		},
	}
	SetProvider(provider)

	enabled, err := provider.IsPairEnabled(provider.Name, "USD/EUR")
	if err != nil {
		t.Fatalf("Error checking if pair is enabled: %v", err)
	}

	if !enabled {
		t.Errorf("Expected pair to be enabled; got disabled")
	}
}
