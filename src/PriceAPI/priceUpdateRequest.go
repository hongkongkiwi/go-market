package PriceAPI

import "fmt"

type PriceUpdateRequest struct {
	Provider  string  `json:"provider"`
	Base      string  `json:"base"`
	Quote     string  `json:"quote"`
	Bid       float64 `json:"bid"`
	BidAmount float64 `json:"bid_amount"`
	Ask       float64 `json:"ask"`
	AskAmount float64 `json:"ask_amount"`
	Timestamp int64   `json:"timestamp"`
}

func (req *PriceUpdateRequest) NewPriceUpdateAsk() *PriceUpdate {
	return &PriceUpdate{
		Provider:  req.Provider,
		Base:      req.Base,
		Quote:     req.Quote,
		Timestamp: req.Timestamp,
		Price:     req.Ask,
		Amount:    req.AskAmount,
	}
}

func (req *PriceUpdateRequest) NewPriceUpdateBid() *PriceUpdate {
	return &PriceUpdate{
		Provider:  req.Provider,
		Base:      req.Base,
		Quote:     req.Quote,
		Timestamp: req.Timestamp,
		Price:     req.Bid,
		Amount:    req.BidAmount,
	}
}

// GetPairName returns the pair name based on the Base and Quote fields of the PriceUpdateRequest.
//
// No parameters.
// Returns a string.
func (p *PriceUpdateRequest) GetPairName() string {
	return fmt.Sprintf("%s/%s", p.Base, p.Quote)
}

func (p *PriceUpdateRequest) GetSpread() float64 {
	return p.Ask - p.Bid
}
