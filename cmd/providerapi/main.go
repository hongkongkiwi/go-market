package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfigAPI"
)

// Constants
const (
	listenAddress = ":8081"
	serverName    = "ProviderAPI"
	dbFile        = "./data/ProviderDB.sqlite"
)

func main() {
	err := ProviderConfig.OpenDB(dbFile)
	if err != nil {
		panic(err)
	}
	defer ProviderConfig.CloseDB()

	if envVar := os.Getenv("PRICE_API_URL_BASE"); envVar != "" {
		ProviderConfigAPI.PriceAPIURLBase = envVar
	}

	// Reshuffle
	ProviderConfig.RandomizeProviders()

	router := SetupRouter()

	// Start the HTTP server
	if err := router.Run(listenAddress); err != nil {
		fmt.Printf("Failed to start %s server: %v\n", serverName, err)
	}
}

func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/ping"),
		gin.Recovery(),
	)

	// Set trusted proxies to fix annoying error
	router.SetTrustedProxies([]string{"127.0.0.1/8"})

	// GET route to retrieve enabled currency pairs for all providers
	router.GET("/providers", ProviderConfigAPI.GetProviders)

	// PUT route to enable/disable currency pairs for a specific provider
	router.PUT("/providers/:providerName", ProviderConfigAPI.SetPairsForProvider)

	// GET route to retrieve enabled currency pairs for a specific provider
	router.GET("/providers/:providerName", ProviderConfigAPI.GetPairsForProvider)

	return router
}
