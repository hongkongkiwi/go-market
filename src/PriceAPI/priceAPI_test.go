package PriceAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
)

func TestProcessPriceUpdate(t *testing.T) {
	tmpFile, tmpFileErr := os.CreateTemp("", "templogfile")
	if tmpFileErr != nil {
		fmt.Println("Error creating temporary file:", tmpFileErr)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up: delete the temporary file when done
	PriceUpdatesLogFile = tmpFile.Name()

	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestProcessPriceUpdate")
	if tempErr != nil {
		t.Errorf("Error creating temporary file: %v", tempErr)
		return
	}
	err := ProviderConfig.OpenDB(tmpDBFileName.Name())
	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer func() {
		ProviderConfig.CloseDB()
		os.Remove(tmpDBFileName.Name())
	}()
	// Create a test PriceUpdateRequest
	update := PriceUpdateRequest{
		Provider:  "TestProvider",
		Base:      "BTC",
		Quote:     "USD",
		Bid:       45000,
		BidAmount: 1,
		Ask:       45500,
		AskAmount: 1,
		Timestamp: 1615299600, // 2021-03-09 00:00:00 UTC
	}
	body, _ := json.Marshal(update)

	// Create a mock HTTP request
	req, err := http.NewRequest("POST", "/price/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a fake Gin context
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	// Call the handler function to process the price update
	ProcessPriceUpdateRequest(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestReCalculateBestPrices(t *testing.T) {
	tmpFile, tmpFileErr := os.CreateTemp("", "templogfile")
	if tmpFileErr != nil {
		fmt.Println("Error creating temporary file:", tmpFileErr)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up: delete the temporary file when done
	PriceUpdatesLogFile = tmpFile.Name()

	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestReCalculateBestPrices")
	if tempErr != nil {
		t.Errorf("Error creating temporary file: %v", tempErr)
		return
	}
	err := ProviderConfig.OpenDB(tmpDBFileName.Name())
	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer func() {
		ProviderConfig.CloseDB()
		os.Remove(tmpDBFileName.Name())
	}()
	ProviderConfig.RandomizeProviders()

	// Create a mock HTTP request
	req, err := http.NewRequest("PUT", "/prices/recalculate", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a fake Gin context
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	// Call the handler function to recalculate best prices
	ReCalculateBestPrices(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}
