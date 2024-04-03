package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hongkongkiwi/chaostheory/src/Helpers"
	"github.com/hongkongkiwi/chaostheory/src/PriceAPI"
	"github.com/hongkongkiwi/chaostheory/src/ProviderConfig"
)

// Constants
const (
	listenAddress = ":8080"
	serverName    = "PriceAPI"
	dbFile        = "./data/ProviderDB.sqlite"
)

func main() {
	// In this case we only need readonly as we don't do any write operations
	// and a limitation of badgerdb is that we can only have one writer
	err := ProviderConfig.OpenDB(dbFile)
	if err != nil {
		panic(err)
	}
	defer ProviderConfig.CloseDB()

	if err := Helpers.CreateDirIfNotExist(filepath.Dir(PriceAPI.PriceUpdatesLogFile)); err != nil {
		panic(err)
	}

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

	// POST route to receive price updates
	router.POST("/prices", PriceAPI.ProcessPriceUpdateRequest)

	// PUT route to recalculate best prices
	router.PUT("/prices/recalculate", PriceAPI.ReCalculateBestPrices)

	// Ping route to check server status
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}
