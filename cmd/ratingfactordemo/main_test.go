package main

import (
	"os"
	"testing"
)

func TestCustomerRatingFactors_GetRatingFactor(t *testing.T) {
	customerRatingFactors := &CustomerRatingFactors{
		RatingFactors: map[string]float64{
			"customer1": 0.8,
			"customer2": 1.0,
			"customer3": 1.2,
		},
	}

	tests := []struct {
		customerID string
		expected   float64
	}{
		{"customer1", 0.8},
		{"customer2", 1.0},
		{"customer3", 1.2},
		{"nonexistent", 1.0}, // Default value
	}

	for _, test := range tests {
		actual := customerRatingFactors.GetRatingFactor(test.customerID)
		if actual != test.expected {
			t.Errorf("GetRatingFactor(%s): expected %.1f, but got %.1f", test.customerID, test.expected, actual)
		}
	}
}

func TestCustomerRatingFactors_GetDisplayedCustomerQuantity(t *testing.T) {
	customerRatingFactors := &CustomerRatingFactors{
		RatingFactors: map[string]float64{
			"customer1": 0.8,
			"customer2": 1.0,
			"customer3": 1.2,
		},
	}

	latestQuantity = 500.0

	tests := []struct {
		customerID string
		expected   float64
	}{
		{"customer1", 500 * 0.8},
		{"customer2", 500 * 1.0},
		{"customer3", 500 * 1.2},
	}

	for _, test := range tests {
		actual := customerRatingFactors.GetDisplayedCustomerQuantity(test.customerID)
		if actual != test.expected {
			t.Errorf("GetDisplayedCustomerQuantity(%s): expected %.2f, but got %.2f", test.customerID, test.expected, actual)
		}
	}
}

func TestCustomerRatingFactors_PlaceOrder(t *testing.T) {
	customerRatingFactors := &CustomerRatingFactors{
		RatingFactors: map[string]float64{
			"customer1": 0.8,
			"customer2": 1.0,
			"customer3": 1.2,
		},
	}

	tests := []struct {
		customerID        string
		displayedQuantity float64
		expectedOutput    string
	}{
		{"customer1", 200.0, "-> Order Placed for Customer customer1 (Inexperienced): Displayed Quantity: 200.00, Adjusted Quantity: 250.00, Follow trade should be made for 50.00\n"},
		{"customer2", 300.0, "-> Order Placed for Customer customer2 (Neutral): Displayed Quantity: 300.00, Adjusted Quantity: 300.00\n"},
		{"customer3", 400.0, "-> Order Placed for Customer customer3 (Experienced): Displayed Quantity: 400.00, Adjusted Quantity: 333.33\n"},
	}

	for _, test := range tests {
		output := captureOutput(func() {
			customerRatingFactors.PlaceOrder(test.customerID, test.displayedQuantity)
		})
		if output != test.expectedOutput {
			t.Errorf("PlaceOrder(%s, %.2f): expected output '%s', but got '%s'", test.customerID, test.displayedQuantity, test.expectedOutput, output)
		}
	}
}

// Helper function to capture fmt.Printf output
func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out := make(chan string)
	go func() {
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		out <- string(buf[:n])
	}()
	os.Stdout = rescueStdout
	return <-out
}
