package ProviderConfigAPI

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"

	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
)

var PriceAPIURLBase = "http://localhost:8080"

type CurrencyPairs struct {
	Base    string `json:"base"`
	Quote   string `json:"quote"`
	Enabled bool   `json:"enabled"`
}

// ProviderPairEnableRequest represents a request to toggle enabled currency pairs for a provider.
type ProviderPairEnableRequest struct {
	Pairs []*CurrencyPairs `json:"pairs"`
}

// setPairsForProvider sets the enabled/disabled currency pairs for a provider.
func SetPairsForProvider(c *gin.Context) {
	providerName := c.Param("providerName")
	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty provider param"})
		return
	}

	// Parse the JSON payload to get the provider and pair status
	var req ProviderPairEnableRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	changedPairs := make(map[string]bool)
	for _, changedPair := range req.Pairs {
		changedPairs[changedPair.Base+"/"+changedPair.Quote] = changedPair.Enabled
	}

	// Update our internal store
	err := ProviderConfig.SetPairsEnabled(providerName, changedPairs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Make a REST client call to /prices/recalculate
	if err := recalculatePrices(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// getPairsForProvider retrieves the current currency pairs and their status for a provider.
func GetPairsForProvider(c *gin.Context) {
	// Extract the provider name from the URL parameter
	providerName := c.Param("providerName")
	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty provider param"})
		return
	}

	provider, err := ProviderConfig.GetProvider(providerName)
	if err != nil {
		log.Printf("Got an error: %e", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pairs := make(map[string]bool)
	if provider != nil {
		pairs = provider.Pairs
	}
	c.JSON(http.StatusOK, pairs)
}

// getProviders retrieves the currency pairs for all providers.
func GetProviders(c *gin.Context) {
	var allProviders map[string]*ProviderConfig.Provider
	var err error
	allProviders, err = ProviderConfig.GetProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the enabled pairs for all providers as JSON response
	c.JSON(http.StatusOK, allProviders)
}

// recalculatePrices makes a PUT to /prices/recalculate
func recalculatePrices() error {
	// Create a new gorequest instance
	request := gorequest.New()

	// Send a PUT request with the scheme and URL specified
	resp, _, errs := request.Put(fmt.Sprintf("%s/prices/recalculate", PriceAPIURLBase)).
		End()

	// Check for errors
	if len(errs) > 0 {
		return fmt.Errorf("request error: %v", errs[0])
	}

	// Check the status code
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("prices/recalculate endpoint not found")
	case http.StatusInternalServerError:
		return fmt.Errorf("internal server error occurred")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
