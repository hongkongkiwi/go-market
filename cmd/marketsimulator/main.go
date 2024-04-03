package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/hongkongkiwi/chaostheory/src/MarketSimulatorClient"
	"github.com/hongkongkiwi/chaostheory/src/MarketSimulatorConfig"
)

func main() {
	config := MarketSimulatorConfig.GenerateDefaultConfig()

	if envVar := os.Getenv("PRICE_API_URL_BASE"); envVar != "" {
		config.APIURL = envVar
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Endlessly loop sending updates
	for {
		// Start simulating provider config updates
		MarketSimulatorClient.StartSimulation(config)

		// Sleep for random time before the next update
		minSleep := config.MinSleep
		maxSleep := config.MaxSleep
		sleepTime := r.Intn(maxSleep-minSleep+1) + minSleep
		time.Sleep(time.Duration(sleepTime) * time.Second) // Adjust sleep time according to update frequency
	}
}
