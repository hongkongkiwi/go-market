package ProviderConfigAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct{}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Create a mock response with status OK
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}, nil
}

func TestSetPairsForProvider(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestSetPairsForProvider")
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

	// Define the pairs struct for test data
	var testPairs = []*CurrencyPairs{
		{Base: "BTC", Quote: "USD", Enabled: true},
		{Base: "ETH", Quote: "USD", Enabled: true},
		{Base: "LTC", Quote: "USD", Enabled: false},
		{Base: "EUR", Quote: "USD", Enabled: true},
		{Base: "XRP", Quote: "USD", Enabled: false},
	}

	// Create test data
	providerName := "DragonFlyExchange"

	// Initialize a new Gin router
	router := gin.Default()

	// Define our routes and handlers
	router.PUT("/providers/:providerName", SetPairsForProvider)

	// Mock the recalculate endpoint (this needs to be callable as it's from another function
	// this is messy but should work
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is PUT and the request path is /prices/recalculate
		if r.Method != http.MethodPut || r.URL.Path != "/prices/recalculate" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		// Respond with "OK"
		w.Write([]byte("OK"))
	}))
	defer mockServer.Close()

	// Create request body using testPairs struct
	reqBody, err := json.Marshal(ProviderPairEnableRequest{Pairs: testPairs})
	if err != nil {
		t.Fatal(err)
	}

	// Create a fake HTTP request with the desired provider name and request body
	req, _ := http.NewRequest("PUT", "/providers/"+providerName, bytes.NewBuffer(reqBody))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Change our recalculate url which gets called as part of the setPairs method
	PriceAPIURLBase = mockServer.URL

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code for setting pairs
	assert.Equal(t, http.StatusOK, rr.Code)
}

func conveertTestPairsToString(testPairs []*CurrencyPairs) string {
	// Initialize a map to store pairs
	pairsMap := make(map[string]bool)

	// Convert testPairs to pairsMap
	for _, pair := range testPairs {
		pairsMap[fmt.Sprintf("%s/%s", pair.Base, pair.Quote)] = pair.Enabled
	}

	// Marshal pairsMap to JSON
	pairsJSON, err := json.Marshal(pairsMap)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(pairsJSON)
}

func TestGetPairsForProvider(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestGetPairsForProvider")
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

	// Define the pairs struct for test data
	var testPairs = []*CurrencyPairs{
		{Base: "BTC", Quote: "USD", Enabled: true},
		{Base: "ETH", Quote: "USD", Enabled: true},
		{Base: "LTC", Quote: "USD", Enabled: false},
		{Base: "EUR", Quote: "USD", Enabled: true},
		{Base: "XRP", Quote: "USD", Enabled: false},
	}

	// Create test data
	providerName := "DragonFlyExchange"
	providerPairs := make(map[string]bool)
	for _, testPair := range testPairs {
		pairStr := fmt.Sprintf("%s/%s", testPair.Base, testPair.Quote)
		providerPairs[pairStr] = testPair.Enabled
	}
	provider := &ProviderConfig.Provider{
		Name:  providerName,
		Pairs: providerPairs,
	}
	ProviderConfig.SetProvider(provider)

	// Initialize a new Gin router
	router := gin.Default()

	// Define our routes and handlers
	router.GET("/providers/:providerName", GetPairsForProvider)

	// Create a fake HTTP request with the desired provider name
	req, _ := http.NewRequest("GET", "/providers/"+providerName, nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the response status code for getting pairs
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body for appropriate pairs
	assert.Equal(t, conveertTestPairsToString(testPairs), rr.Body.String()) // Assuming the response body is JSON-encoded pairs
}

func TestGetProviders(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestGetProviders")
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

	// Create test data
	expectedProviders, err := ProviderConfig.RandomizeProviders()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/providers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a fake Gin context
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	// Call the handler function to get providers
	GetProviders(c)

	// Check the response status code for getting providers
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var actualProviders map[string]*ProviderConfig.Provider
	if err := json.Unmarshal(rr.Body.Bytes(), &actualProviders); err != nil {
		t.Fatal(err)
	}

	// Compare the retrieved providers with the expected ones
	assert.Equal(t, expectedProviders, actualProviders)
}
