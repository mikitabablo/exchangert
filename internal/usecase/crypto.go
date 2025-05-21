package usecase

import (
	"errors"
	"math"

	"github.com/mikitabablo/exchangert/internal/domain"
)

type CryptoUsecase struct {
	provider domain.CryptoRatesProvider
}

func NewCryptoUsecase(provider domain.CryptoRatesProvider) *CryptoUsecase {
	return &CryptoUsecase{provider: provider}
}

func (c *CryptoUsecase) Exchange(from, to string, amount float64) (domain.CryptoExchangeResult, error) {
	if amount <= 0 {
		return domain.CryptoExchangeResult{}, errors.New("invalid amount")
	}

	rates := c.provider.GetCryptoRates()
	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 || !ok2 {
		return domain.CryptoExchangeResult{}, errors.New("unsupported currency")
	}

	usd := amount * fromRate.RateUSD
	converted := usd / toRate.RateUSD

	precision := math.Pow(10, float64(toRate.DecimalPlaces))
	converted = math.Round(converted*precision) / precision

	return domain.CryptoExchangeResult{
		From:   from,
		To:     to,
		Amount: converted,
	}, nil
}
