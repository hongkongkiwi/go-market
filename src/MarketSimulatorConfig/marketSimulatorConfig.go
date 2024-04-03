package MarketSimulatorConfig

import (
	"fmt"
	"strings"

	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
)

// CurrencyPair represents a currency pair with base and quote currencies.
type CurrencyPair [2]string

// String returns the string representation of the currency pair.
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s", cp[0], cp[1])
}

var DbFile = "./data/ProviderDB.sqlite"

// SimulatorConfig represents the configuration for the price simulator.
type SimulatorConfig struct {
	// Default API URL for price updates
	APIURL string `yaml:"api_url"`
	// Minimum spread allowed between bid and ask prices
	MinSpread float64 `yaml:"min_spread"`
	// Minimum sleep time between updates (in seconds)
	MinSleep int `yaml:"min_sleep"`
	// Maximum sleep time between updates (in seconds)
	MaxSleep int `yaml:"max_sleep"`
	// Default providers with associated currency pairs
	Providers map[string][]*CurrencyPair `yaml:"providers"`
	// Default initial price for assets
	InitialPrice float64 `yaml:"initial_price"`
	// Default initial quantity for assets
	InitialQuantity float64 `yaml:"initial_quantity"`
	// Default factor for price change in each update
	PriceChangeFactor    float64 `yaml:"price_change_factor"`
	QuantityChangeFactor float64 `yaml:"quantity_change_factor"`
	AllowArbitrage       bool    `yaml:"allow_arbitrage"`
}

// generateDefaultConfig generates default configuration data and writes it to a file.
func GenerateDefaultConfig() *SimulatorConfig {
	ProviderConfig.OpenDB(DbFile)
	defer ProviderConfig.CloseDB()
	providers, err := ProviderConfig.GetProviders()
	if err != nil {
		ProviderConfig.RandomizeProviders()
	}
	if len(providers) == 0 {
		ProviderConfig.RandomizeProviders()
	}

	// We just pull some providers from our db list
	// this saves redefining them here, probably we
	// would use a json config file or something rather
	// than pulling from the real providers db but
	// I am a bit lazy in this app
	providersConfig := make(map[string][]*CurrencyPair)
	for providerName, provider := range providers {
		for pairName := range provider.Pairs {
			currencyPair := &CurrencyPair{}
			currencyPair[0] = strings.Split(pairName, "/")[0]
			currencyPair[1] = strings.Split(pairName, "/")[1]
			providersConfig[providerName] = append(providersConfig[providerName], currencyPair)
		}
	}

	// Define default configuration data
	defaultConfig := &SimulatorConfig{
		// Default API URL for price updates
		APIURL: "http://localhost:8080",
		// Minimum spread allowed between bid and ask prices
		MinSpread: 0.1,
		// Minimum sleep time between updates (in seconds)
		MinSleep: 1,
		// Maximum sleep time between updates (in seconds)
		MaxSleep: 5,
		// Default providers with associated currency pairs
		Providers: providersConfig,
		// Default initial price for assets
		InitialPrice:    100.0,
		InitialQuantity: 1000.0,
		// Default factor for price change in each update
		PriceChangeFactor:    10.0,
		QuantityChangeFactor: 100.0,
		AllowArbitrage:       false,
	}

	return defaultConfig
}
