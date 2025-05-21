package crypto

import "github.com/mikitabablo/exchangert/internal/domain"

type StaticDataProvider struct{}

func NewStaticDataProvider() *StaticDataProvider {
	return &StaticDataProvider{}
}

func (h *StaticDataProvider) GetCryptoRates() map[string]domain.CryptoCurrency {
	return map[string]domain.CryptoCurrency{
		"BEER":  {Name: "BEER", DecimalPlaces: 18, RateUSD: 0.00002461},
		"FLOKI": {Name: "FLOKI", DecimalPlaces: 18, RateUSD: 0.0001428},
		"GATE":  {Name: "GATE", DecimalPlaces: 18, RateUSD: 6.87},
		"USDT":  {Name: "USDT", DecimalPlaces: 6, RateUSD: 0.999},
		"WBTC":  {Name: "WBTC", DecimalPlaces: 8, RateUSD: 57037.22},
	}
}
