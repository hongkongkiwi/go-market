package ProviderConfig

import "encoding/json"

// Provider represents a data provider.
type Provider struct {
	Name  string          `json:"provider"`
	Pairs map[string]bool `json:"pairs"` // Mapping of currency pairs strings to enabled/disabled
}

// SerializeToBytes serializes the Provider struct to bytes using JSON encoding.
func (p *Provider) SerializeToBytes() ([]byte, error) {
	return json.Marshal(p)
}

// DeserializeFromBytes deserializes the Provider struct from bytes using JSON encoding.
func DeserializeProviderFromBytes(data []byte) (*Provider, error) {
	var provider Provider
	err := json.Unmarshal(data, &provider)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetPairs retrieves the currency pairs for a given provider.
func (provider *Provider) GetPairs(providerName string) (map[string]bool, error) {
	provider, err := GetProvider(providerName)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return make(map[string]bool, 0), nil
	}
	return provider.Pairs, nil
}

func (provider *Provider) IsPairEnabled(providerName string, pair string) (bool, error) {
	provider, err := GetProvider(providerName)
	if err != nil {
		return false, err
	}
	if provider == nil {
		return false, nil
	}
	return provider.Pairs[pair], nil
}
