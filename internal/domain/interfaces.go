package domain

import "context"

type FiatRatesProvider interface {
    GetRates(base string, symbols []string) (*RatesResponse, error)
}

type FiatUsecase interface {
	GetRates(ctx context.Context, currencies []string) ([]FiatRate, error)
}

type CryptoRatesProvider interface {
	GetCryptoRates() map[string]CryptoCurrency
}

type CryptoUsecase interface {
	Exchange(from, to string, amount float64) (CryptoExchangeResult, error)
}
