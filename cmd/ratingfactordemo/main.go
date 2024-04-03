package main

import (
	"fmt"
	"sync"
)

var latestQuantity = 500.0

type CustomerRatingFactors struct {
	RatingFactors map[string]float64
	mu            sync.RWMutex
}

type CustomerType string

const (
	Experienced   CustomerType = "Experienced"
	Neutral       CustomerType = "Neutral"
	Inexperienced CustomerType = "Inexperienced"
)

func (c *CustomerRatingFactors) GetRatingFactor(customerID string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if factor, ok := c.RatingFactors[customerID]; ok {
		return factor
	}
	// Return a default rating factor if non exists
	return 1.0
}

func (c *CustomerRatingFactors) GetDisplayedCustomerQuantity(customerID string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	factor := c.GetRatingFactor(customerID)
	return latestQuantity * factor
}

func (c *CustomerRatingFactors) GetActualBackendQuantity(customerID string, displayedQuantity float64) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	factor := c.GetRatingFactor(customerID)
	return displayedQuantity / factor
}

func (c *CustomerRatingFactors) PlaceOrder(customerID string, displayedQuantity float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	factor := c.GetRatingFactor(customerID)
	adjustedQuantity := c.GetActualBackendQuantity(customerID, displayedQuantity)
	customerType := c.GetCustomerType(factor)
	fmt.Printf("-> Order Placed for Customer %s (%s): Displayed Quantity: %.2f, Adjusted Quantity: %.2f", customerID, customerType, displayedQuantity, adjustedQuantity)
	if factor < 1.0 {
		followTradeQuantity := adjustedQuantity - displayedQuantity
		fmt.Printf(", Follow trade should be made for %.2f\n", followTradeQuantity)
	} else {
		fmt.Println() // Just print a newline if no follow trade is needed
	}
}

// This is not necessary for the problem, just
// makes it nice for logging
func (c *CustomerRatingFactors) GetCustomerType(factor float64) CustomerType {
	switch {
	case factor < 1.0:
		return Inexperienced
	case factor == 1.0:
		return Neutral
	default:
		return Experienced
	}
}

func ShowCustomerDisplayedQuantity(customer string, customerRatingFactors *CustomerRatingFactors) {
	displayedQuantity := customerRatingFactors.GetDisplayedCustomerQuantity(customer)
	ratingFactor := customerRatingFactors.GetRatingFactor(customer)
	customerType := customerRatingFactors.GetCustomerType(ratingFactor)
	fmt.Printf("-> Customer %s: Type: %s, Displayed Quantity: %.2f, Rating Factor: %.1f\n", customer, customerType, displayedQuantity, ratingFactor)
}

func main() {
	fmt.Println("Running RatingFactorDemo")

	customerRatingFactors := &CustomerRatingFactors{
		RatingFactors: map[string]float64{
			"customer1": 0.8,
			"customer2": 1.0,
			"customer3": 1.2,
		},
	}

	fmt.Println("Requesting the liquidity available to customers:")
	// Just for demo purposes loop through all customers and request the best quantity
	for customer := range customerRatingFactors.RatingFactors {
		ShowCustomerDisplayedQuantity(customer, customerRatingFactors)
	}

	// Create a bunch of fake orders, this is truncated just to show quantity adjustment
	orders := []struct {
		CustomerID string
		Quantity   float64
	}{
		{CustomerID: "customer1", Quantity: 200.0},
		{CustomerID: "customer2", Quantity: 300.0},
		{CustomerID: "customer3", Quantity: 400.0},
	}

	fmt.Println("Placing example orders for customers:")
	for _, order := range orders {
		customerRatingFactors.PlaceOrder(order.CustomerID, order.Quantity)
	}
}
