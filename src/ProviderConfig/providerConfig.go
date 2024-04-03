package ProviderConfig

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDB(dbFile string) error {
	dbDir := path.Dir(dbFile)
	// Check if the directory exists
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		// Directory does not exist, so create it
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return err
		}
	}

	// Open or create the database file
	sqliteDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	// Create the table if it doesn't exist
	_, err = sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS providers (
			name TEXT PRIMARY KEY,
			"pairs" TEXT
		);`)
	if err != nil {
		sqliteDB.Close()
		return err
	}

	// Set the global database variable
	db = sqliteDB

	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func SetProvider(provider *Provider) error {
	// Null check for provider
	if provider == nil {
		return fmt.Errorf("provider is nil")
	}

	// If pairs are nil, create an empty map
	if provider.Pairs == nil {
		provider.Pairs = make(map[string]bool)
	}

	// Add the provider to the database
	stmt, err := db.Prepare("REPLACE INTO providers (name, pairs) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pairsJSON, err := json.Marshal(provider.Pairs)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(provider.Name, string(pairsJSON))
	if err != nil {
		return err
	}

	return nil
}

func GetProviders() (map[string]*Provider, error) {
	providers := make(map[string]*Provider)

	rows, err := db.Query("SELECT name, pairs FROM providers")
	if err != nil {
		return providers, err
	}
	defer rows.Close()

	for rows.Next() {
		var providerName string
		var pairsJSON string

		err := rows.Scan(&providerName, &pairsJSON)
		if err != nil {
			return providers, err
		}

		provider := &Provider{
			Name: providerName,
		}

		err = json.Unmarshal([]byte(pairsJSON), &provider.Pairs)
		if err != nil {
			return providers, err
		}

		providers[providerName] = provider
	}

	return providers, nil
}

func GetProvider(providerName string) (*Provider, error) {
	var provider *Provider

	row := db.QueryRow("SELECT pairs FROM providers WHERE name = ?", providerName)

	var pairsJSON string
	err := row.Scan(&pairsJSON)
	if err == sql.ErrNoRows {
		// This is not an error, just no provider
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	provider = &Provider{
		Name: providerName,
	}

	err = json.Unmarshal([]byte(pairsJSON), &provider.Pairs)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// GetProviderPairEnabled retrieves the enabled status of a specific currency pair for a given provider.
func GetProviderPairEnabled(providerName string, pairName string) (bool, error) {
	// Retrieve the provider
	provider, err := GetProvider(providerName)
	if err != nil {
		return false, err
	}
	if provider == nil {
		// Provider not found, return false
		return false, nil
	}
	// Retrieve the enabled status of the pair
	enabled, ok := provider.Pairs[pairName]
	if !ok {
		// Pair not found, return false
		return false, nil
	}
	return enabled, nil
}

// SetPairEnabled enables or disables a specific currency pair for a given provider.
func SetPairEnabled(providerName string, pair string, enabled bool) error {
	provider, err := GetProvider(providerName)
	if err != nil {
		return err
	}
	if provider == nil {
		provider = &Provider{
			Name:  providerName,
			Pairs: make(map[string]bool, 0),
		}
	}
	provider.Pairs[pair] = enabled
	return SetProvider(provider)
}

func SetPairsEnabled(providerName string, pairsEnabled map[string]bool) error {
	fmt.Printf("Setting pairs for %s: %v\n", providerName, pairsEnabled)
	provider, err := GetProvider(providerName)
	if err != nil {
		return err
	}
	if provider == nil {
		provider = &Provider{
			Name:  providerName,
			Pairs: make(map[string]bool, 0),
		}
	}
	for pair, enabled := range pairsEnabled {
		provider.Pairs[pair] = enabled
	}
	return SetProvider(provider)
}

// generateRandomProviderSettings generates random enabled/disabled states
// for each provider and pair listed below
func RandomizeProviders() (map[string]*Provider, error) {
	// Create a new source using the current Unix timestamp
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator using the source
	r := rand.New(source)

	// Define our default set of providers here
	defaultProviders := map[string][]string{
		"DragonFlyExchange":      {"BTC/USD", "ETH/USD", "LTC/USD", "EUR/USD", "XRP/USD", "BTC/EUR", "ETH/EUR", "LTC/EUR", "EUR/GBP", "XRP/GBP"},
		"MoonlightExchange":      {"BTC/USD", "ETH/USD", "LTC/USD", "XRP/USD", "BCH/USD", "BTC/GBP", "ETH/GBP", "LTC/GBP", "EUR/GBP", "BCH/GBP"},
		"StellarHorizonExchange": {"BTC/USD", "ETH/USD", "LTC/USD", "XRP/USD", "BCH/USD", "BTC/JPY", "ETH/JPY", "LTC/JPY", "XRP/JPY", "BCH/JPY"},
		"AuroraExchange":         {"BTC/USD", "ETH/USD", "XRP/USD", "EUR/USD", "BCH/USD", "BTC/JPY", "ETH/JPY", "XRP/JPY", "EUR/JPY", "BCH/JPY"},
		"GalacticExchange":       {"BTC/USD", "LTC/USD", "EUR/USD", "XRP/USD", "BCH/USD", "BTC/AUD", "LTC/AUD", "EUR/AUD", "XRP/AUD", "BCH/AUD"},
		"GoldenDragonExchange":   {"BTC/USD", "ETH/USD", "XRP/USD", "EUR/USD", "BCH/USD", "BTC/GBP", "LTC/GBP", "EUR/GBP", "XRP/GBP", "BCH/GBP", "BTC/AUD", "ETH/AUD"},
	}

	providers := make(map[string]*Provider)
	// Assign enabled or disabled status randomly for each pair
	for providerName := range defaultProviders {
		provider := &Provider{
			Name:  providerName,
			Pairs: make(map[string]bool),
		}
		for _, pairName := range defaultProviders[providerName] {
			provider.Pairs[pairName] = r.Intn(3) > 0 // Give a 2/3 change to be enabled to illustrate our app better
		}
		// This is not the best handling here since we would finish the loop
		// immediately upon error but better handling is not really necessary here
		if err := SetProvider(provider); err != nil {
			return nil, err
		}
		providers[providerName] = provider
		fmt.Printf("New Randomized Provider %s: %v\n", providerName, provider.Pairs)
	}
	return providers, nil
}
