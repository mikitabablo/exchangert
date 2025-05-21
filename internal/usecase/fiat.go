package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/mikitabablo/exchangert/internal/domain"
)

type RatesUsecase struct {
	provider domain.FiatRatesProvider
}

func NewRatesUsecase(p domain.FiatRatesProvider) *RatesUsecase {
	return &RatesUsecase{provider: p}
}

func (r *RatesUsecase) GetRates(ctx context.Context, currencies []string) ([]domain.FiatRate, error) {
	if len(currencies) < 2 {
		return nil, errors.New("at least 2 currencies required")
	}

	resp, err := r.provider.GetRates(domain.BaseCurrency, currencies)
	if err != nil {
		return nil, err
	}

	rates := resp.Rates
	base := resp.Base
	rates[base] = 1.0

	filteredRates := make(map[string]float64)
	for _, currency := range currencies {
		rate, ok := rates[currency]
		if !ok {
			return nil, fmt.Errorf("currency %s not found in rates", currency)
		}
		filteredRates[currency] = rate
	}

	var result []domain.FiatRate
	for _, from := range currencies {
		rateFrom := filteredRates[from]
		for _, to := range currencies {
			if from == to {
				continue
			}
			rateTo := filteredRates[to]
			converted := rateTo / rateFrom
			result = append(result, domain.FiatRate{
				From: from,
				To:   to,
				Rate: converted,
			})
		}
	}

	return result, nil
}
