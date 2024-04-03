package MarketSimulatorClient

import (
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hongkongkiwi/chaostheory/src/MarketSimulatorConfig"
	"github.com/hongkongkiwi/chaostheory/src/PriceAPI"
	"github.com/parnurzeal/gorequest"
)

type PriceHistory struct {
	BidPrice  float64
	AskPrice  float64
	BidAmount float64
	AskAmount float64
}

var priceHistories = make(map[string]map[string]*PriceHistory)

func StartSimulation(config *MarketSimulatorConfig.SimulatorConfig) {
	if config == nil {
		log.Println("Error: nil configuration provided to StartSimulation")
		return
	}

	// Create a new random number generator with a specific seed value
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Start simulating provider config updates
	for provider, currencyPairs := range config.Providers {
		// fmt.Printf("Starting simulation for provider: %s\n", provider)
		if priceHistories[provider] == nil {
			priceHistories[provider] = make(map[string]*PriceHistory)
		}
		for _, currencyPair := range currencyPairs {
			// fmt.Printf("Starting simulation for currency pair: %s\n", currencyPair.String())
			// Generate random ask and bid prices
			ask, bid := generateNewPricesForPair(provider, currencyPair, r, config)

			// Generate random ask and bid amounts based on previous values
			askAmount, bidAmount := generateNewAmountsForPair(provider, currencyPair, r, config)

			// Log the ask and bid prices
			log.Printf("%s (%s) Ask price: %.2f amount: %.2f, Bid price: %.2f amount: %.2f", provider, currencyPair.String(), ask, askAmount, bid, bidAmount)

			// Prepare JSON payload
			payload := &PriceAPI.PriceUpdateRequest{
				Provider:  provider,
				Base:      currencyPair[0],
				Quote:     currencyPair[1],
				Bid:       bid,
				BidAmount: bidAmount,
				Ask:       ask,
				AskAmount: askAmount,
				Timestamp: time.Now().UnixNano(),
			}

			// Send POST request to price store API endpoint using gorequest
			resp, _, errs := gorequest.New().Post(config.APIURL + "/prices").
				Send(payload).
				Type("json").
				End()

			// Check for errors
			if len(errs) > 0 {
				for _, err := range errs {
					log.Printf("Error sending POST request: %v", err)
				}
				continue
			}

			// Check response status code
			if resp.StatusCode == 200 {
				log.Printf("Sent update for provider %s and currency pair %s", provider, currencyPair.String())
			} else {
				body, _ := io.ReadAll(io.LimitReader(resp.Body, int64(os.Getpagesize())))
				log.Printf("Unexpected status code: %d. Response body: %s", resp.StatusCode, string(body))
			}
			resp.Body.Close()
		}
	}
}

func generateNewPricesForPair(provider string, currencyPair *MarketSimulatorConfig.CurrencyPair, r *rand.Rand, config *MarketSimulatorConfig.SimulatorConfig) (float64, float64) {
	// Get the price history for the given provider and currency pair
	pairName := currencyPair.String()
	priceHistory := priceHistories[provider][pairName]
	if priceHistory == nil {
		// If there's no price history, create a new PriceHistory entry with the initial price
		priceHistories[provider][currencyPair.String()] = &PriceHistory{
			AskPrice: config.InitialPrice,
			BidPrice: config.InitialPrice,
		}
		priceHistory = priceHistories[provider][pairName]
	}

	// Grab the last ask price
	lastAskPrice := priceHistory.AskPrice
	// If generating an ask price, increase or decrease the price by a random factor
	newAskPrice := lastAskPrice + (r.Float64()*2-1)*config.PriceChangeFactor
	// Ensure new ask price is above 0.01
	if newAskPrice < 0.01 {
		newAskPrice = 0.01 + (r.Float64()*(config.PriceChangeFactor-0.01) + 0.01)
	}
	if !config.AllowArbitrage {
		// Ensure new ask price is above last bid price so we dont generate arbitrage
		newAskPrice = math.Max(newAskPrice, priceHistory.BidPrice+0.01)
	}

	// Grab the last bid price
	lastBidPrice := priceHistory.BidPrice
	// If generating a bid price, increase or decrease the price by a random factor
	newBidPrice := lastBidPrice + (r.Float64()*2-1)*config.PriceChangeFactor
	// Ensure new bid price is above 0.01
	if newBidPrice < 0.01 {
		newBidPrice = 0.01 + (r.Float64()*(config.PriceChangeFactor-0.01) + 0.01)
	}
	if !config.AllowArbitrage {
		// Ensure new bid price is below last ask price so we dont generate arbitrage opportunity
		newBidPrice = math.Min(newBidPrice, priceHistory.AskPrice-0.01)
	}

	return math.Round(newAskPrice*100) / 100, math.Round(newBidPrice*100) / 100
}

func generateNewAmountsForPair(provider string, currencyPair *MarketSimulatorConfig.CurrencyPair, r *rand.Rand, config *MarketSimulatorConfig.SimulatorConfig) (float64, float64) {
	// Get the price history for the given provider and currency pair
	pairName := currencyPair.String()
	priceHistory := priceHistories[provider][pairName]
	if priceHistory == nil {
		// If there's no price history, create a new PriceHistory entry with the initial price
		priceHistories[provider][currencyPair.String()] = &PriceHistory{
			AskAmount: config.InitialQuantity,
			BidAmount: config.InitialQuantity,
		}
		priceHistory = priceHistories[provider][pairName]
	}

	// Grab the last ask amount
	lastAskAmount := priceHistory.AskAmount
	// If generating an ask price, increase or decrease the price by a random factor
	newAskAmount := lastAskAmount + (r.Float64()*2-1)*config.QuantityChangeFactor
	// Ensure new ask amount is above 0
	if newAskAmount < 0 {
		newAskAmount = 0.01 + (r.Float64()*config.QuantityChangeFactor + 0.01)
	}

	// Grab the last bid amount
	lastBidAmount := priceHistory.BidAmount
	// If generating a bid price, increase or decrease the price by a random factor
	newBidAmount := lastBidAmount + (r.Float64()*2-1)*config.QuantityChangeFactor
	// Ensure new bid amount is above 0
	if newBidAmount < 0 {
		newBidAmount = 0.01 + (r.Float64()*config.QuantityChangeFactor + 0.01)
	}

	return math.Round(newAskAmount*100) / 100, math.Round(newBidAmount*100) / 100
}
