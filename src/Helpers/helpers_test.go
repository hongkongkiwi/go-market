package Helpers

import (
	"os"
	"testing"
)

func TestAppendToFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_append_*.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Remove the temporary file when done

	// Test data
	testData := "This is a test line\n"

	// Test appending to the temporary file
	if err := AppendToFile(tmpFile.Name(), testData); err != nil {
		t.Fatalf("Error appending to file: %v", err)
	}

	// Read the content of the temporary file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Check if the content matches the test data
	if string(content) != testData {
		t.Errorf("Expected content %q; got %q", testData, string(content))
	}

	// Test appending to the file again
	additionalData := "This is additional test line\n"
	if err := AppendToFile(tmpFile.Name(), additionalData); err != nil {
		t.Fatalf("Error appending to file: %v", err)
	}

	// Read the updated content of the temporary file
	updatedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Check if the updated content includes the additional data
	expectedUpdatedContent := testData + additionalData
	if string(updatedContent) != expectedUpdatedContent {
		t.Errorf("Expected updated content %q; got %q", expectedUpdatedContent, string(updatedContent))
	}
}
