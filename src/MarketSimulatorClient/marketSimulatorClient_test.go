package MarketSimulatorClient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/hongkongkiwi/chaostheory/src/MarketSimulatorConfig"
	"github.com/stretchr/testify/assert"
)

func TestStartSimulation(t *testing.T) {
	tmpDBFileName, tempErr := Helpers.CreateTempFile("TestStartSimulation")
	if tempErr != nil {
		t.Errorf("Error creating temporary file: %v", tempErr)
		return
	}
	MarketSimulatorConfig.DbFile = tmpDBFileName.Name()

	// Mock the API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Define a sample configuration for testing
	config := &MarketSimulatorConfig.SimulatorConfig{
		Providers: map[string][]*MarketSimulatorConfig.CurrencyPair{
			"Provider1": {{"BTC", "USD"}, {"ETH", "USD"}},
			"Provider2": {{"BTC", "EUR"}, {"ETH", "EUR"}},
		},
		APIURL:            server.URL,
		MinSleep:          1,
		MaxSleep:          5,
		PriceChangeFactor: 0.01,
	}

	// Test simulation with empty configuration
	emptyConfig := &MarketSimulatorConfig.SimulatorConfig{}
	priceHistories = make(map[string]map[string]*PriceHistory)
	StartSimulation(emptyConfig)
	assert.Empty(t, priceHistories, "priceHistories should remain empty with empty configuration")

	// Test simulation with nil configuration
	priceHistories = make(map[string]map[string]*PriceHistory)
	StartSimulation(nil)
	assert.Empty(t, priceHistories, "priceHistories should remain empty with nil configuration")

	// Test simulation with valid configuration
	priceHistories = make(map[string]map[string]*PriceHistory)
	StartSimulation(config)

	// Check price deviations
	checkPriceDeviation(t, config)

	// Check amount deviations
	checkAmountDeviation(t, config)
}

func checkPriceDeviation(t *testing.T, config *MarketSimulatorConfig.SimulatorConfig) {
	for provider, pairs := range config.Providers {
		for _, pair := range pairs {
			// Get previous prices from history
			prev := priceHistories[provider][pair.String()]
			if prev == nil {
				continue
			}

			// Calculate allowed price deviation
			allowedDeviation := config.PriceChangeFactor

			// Check bid price deviation
			assert.InDelta(t, prev.BidPrice, priceHistories[provider][pair.String()].BidPrice, allowedDeviation, "Bid price deviation exceeds allowed margin")

			// Check ask price deviation
			assert.InDelta(t, prev.AskPrice, priceHistories[provider][pair.String()].AskPrice, allowedDeviation, "Ask price deviation exceeds allowed margin")
		}
	}
}

func checkAmountDeviation(t *testing.T, config *MarketSimulatorConfig.SimulatorConfig) {
	for provider, pairs := range config.Providers {
		for _, pair := range pairs {
			// Get previous amounts from history
			prev := priceHistories[provider][pair.String()]
			if prev == nil {
				continue
			}

			// Calculate allowed amount deviation
			allowedDeviation := config.PriceChangeFactor

			// Check bid amount deviation
			assert.InDelta(t, prev.BidAmount, priceHistories[provider][pair.String()].BidAmount, allowedDeviation, "Bid amount deviation exceeds allowed margin")

			// Check ask amount deviation
			assert.InDelta(t, prev.AskAmount, priceHistories[provider][pair.String()].AskAmount, allowedDeviation, "Ask amount deviation exceeds allowed margin")
		}
	}
}
