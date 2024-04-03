package PriceAPI

import "fmt"

type PriceUpdate struct {
	Provider  string
	Base      string
	Quote     string
	Price     float64
	Amount    float64
	Timestamp int64
}

// GetPairName returns the pair name based on the Base and Quote fields of the PriceUpdateRequest.
//
// No parameters.
// Returns a string.
func (p *PriceUpdate) GetPairName() string {
	return fmt.Sprintf("%s/%s", p.Base, p.Quote)
}
