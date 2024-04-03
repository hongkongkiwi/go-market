package PriceAPI

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
)

var PriceUpdatesLogFile = "./logs/best_prices.log"

type PriceUpdateType = string

var (
	bestBidStore            = make(map[string]*PriceUpdate)
	bestAskStore            = make(map[string]*PriceUpdate)
	providerLastUpdateStore = make(map[string]map[string]*PriceUpdateRequest)
	// Many readers, one writer
	mu sync.RWMutex
)

func GetBestBidPrice(pairName string) *PriceUpdate {
	mu.RLock()
	defer mu.RUnlock()
	return bestBidStore[pairName]
}

func GetBestAskPrice(pairName string) *PriceUpdate {
	mu.RLock()
	defer mu.RUnlock()
	return bestAskStore[pairName]
}

func updateBestBidPrice(newPrice *PriceUpdate) {
	mu.Lock()
	defer mu.Unlock()
	// Should handle this better
	if newPrice == nil {
		return
	}
	bestBidStore[newPrice.GetPairName()] = newPrice
}

func updateBestAskPrice(newPrice *PriceUpdate) {
	mu.Lock()
	defer mu.Unlock()
	// Should handle this better
	if newPrice == nil {
		return
	}
	bestAskStore[newPrice.GetPairName()] = newPrice
}

func saveProviderUpdateRequest(update *PriceUpdateRequest) {
	mu.Lock()
	defer mu.Unlock()
	// Should handle this better
	if update == nil {
		return
	}
	if providerLastUpdateStore[update.Provider] == nil {
		providerLastUpdateStore[update.Provider] = make(map[string]*PriceUpdateRequest)
	}
	providerLastUpdateStore[update.Provider][update.GetPairName()] = update
}

func getLastPriceUpdateRequests(providerName string) map[string]*PriceUpdateRequest {
	mu.RLock()
	defer mu.RUnlock()
	if providerLastUpdateStore[providerName] != nil {
		return make(map[string]*PriceUpdateRequest)
	}
	return providerLastUpdateStore[providerName]
}

func getProviderList() []string {
	mu.RLock()
	defer mu.RUnlock()
	providerNames := make([]string, 0, len(providerLastUpdateStore))
	for providerName := range providerLastUpdateStore {
		providerNames = append(providerNames, providerName)
	}
	return providerNames
}

// ProcessPriceUpdate handles updating the best bid and ask prices.
func ProcessPriceUpdateRequest(c *gin.Context) {
	// Populate our PriceUpdateRequest from received JSON
	var updatePriceReq PriceUpdateRequest
	if err := c.BindJSON(&updatePriceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatePriceReq.Provider == "" || updatePriceReq.Base == "" || updatePriceReq.Quote == "" {
		fmt.Printf("Missing provider, base, or quote fields in PriceUpdateRequest.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing provider, base, or quote fields."})
		return
	}

	// Calculate spread between bid and ask prices
	if updatePriceReq.GetSpread() < 0 {
		// Arbitrage opportunity detected reject update
		fmt.Printf("Arbitrage opportunity detected for update on provider %s, dropping PriceUpdate\n", updatePriceReq.Provider)
		c.JSON(http.StatusBadRequest, "Arbitrage opportunity detected. Dropping PriceUpdate")
		return
	}

	// We can send a result back at this point
	c.Status(http.StatusOK)

	// I am unsure if this needs to be in another thread
	// depends on how Gin works it's contexts but it
	// should be safe to do so
	go func() {
		pairName := updatePriceReq.GetPairName()
		// Save this update so we can use it for recalculation later
		saveProviderUpdateRequest(&updatePriceReq)

		isEnabled, _ := ProviderConfig.GetProviderPairEnabled(updatePriceReq.Provider, pairName)
		// Only update the best price if this provider is enabled
		if isEnabled {
			bestBidPrice := GetBestBidPrice(pairName)
			// Update the best bid and ask prices based on whether this new price is better than the last
			if bestBidPrice == nil || updatePriceReq.Bid > bestBidPrice.Price {
				bidPriceUpdate := updatePriceReq.NewPriceUpdateBid()
				updateBestBidPrice(bidPriceUpdate)
				emitPriceUpdate(bidPriceUpdate, "Bid")
			}

			bestAskPrice := GetBestAskPrice(pairName)
			if bestAskPrice == nil || updatePriceReq.Ask > bestAskPrice.Price {
				askPriceUpdate := updatePriceReq.NewPriceUpdateAsk()
				updateBestBidPrice(askPriceUpdate)
				emitPriceUpdate(askPriceUpdate, "Ask")
			}
		} else {
			// We only log if the provider is enabled
			fmt.Printf("Provider %s is disabled, not updating price for %s\n", updatePriceReq.Provider, pairName)
		}
	}()
}

// recalculatePriceUpdates chooses the best bid and ask prices based on all enabled
// price updates generally this is called when a provider is enabled or disabled
// as it's a bit more expensive than simply checking the previous best price
func ReCalculateBestPrices(c *gin.Context) {
	// We can return a result quickly
	c.Status(http.StatusOK)

	// I am unsure if this needs to be in another thread
	// depends on how Gin works it's contexts but it
	// should be safe to do so
	go func() {
		newAskUpdates := make(map[string]*PriceUpdate)
		newBidUpdates := make(map[string]*PriceUpdate)

		// Loop through all providers we have recevied prices from
		for _, providerName := range getProviderList() {
			// Loop through all currency pairs attached to that provider
			for pairName, updatePriceReq := range getLastPriceUpdateRequests(providerName) {
				if newAskUpdates[pairName] == nil || updatePriceReq.Ask > newAskUpdates[pairName].Price {
					newAskUpdates[pairName] = updatePriceReq.NewPriceUpdateAsk()
				}

				if newBidUpdates[pairName] == nil || updatePriceReq.Bid > newBidUpdates[pairName].Price {
					newBidUpdates[pairName] = updatePriceReq.NewPriceUpdateBid()
				}
			}
		}

		// Update all new ask updates
		for pairName, update := range newAskUpdates {
			// Check again that it's better (a better update may have come in)
			if update.Price > GetBestAskPrice(pairName).Price {
				updateBestAskPrice(update)
				emitPriceUpdate(update, "Ask")
			}
		}

		// Update all new bid updates
		for pairName, update := range newBidUpdates {
			// Check again that it's better (a better update may have come in)
			if update.Price > GetBestBidPrice(pairName).Price {
				updateBestBidPrice(update)
				emitPriceUpdate(update, "Bid")
			}
		}
	}()
}

// emitPriceUpdateUpdate is called when we have a new best price update
// to communicate. In this case we just append to a log file.
func emitPriceUpdate(update *PriceUpdate, updateType string) {
	var logEntry string
	if update != nil {
		logEntry = fmt.Sprintf("%s - %s - %.2f - %.2f - %s\n", updateType, update.Provider, update.Price, update.Amount, time.Unix(update.Timestamp/1000, 0))
	} else {
		logEntry = fmt.Sprintf("%s - No best price available\n", updateType)
	}
	fmt.Println(logEntry)
	// Append to our log file
	if err := Helpers.AppendToFile(PriceUpdatesLogFile, logEntry); err != nil {
		fmt.Println("Error appending to log file:", err)
	}
}
