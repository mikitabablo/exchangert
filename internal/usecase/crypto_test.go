package usecase_test

import (
	"testing"

	"github.com/mikitabablo/exchangert/internal/client/crypto"
	"github.com/mikitabablo/exchangert/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestExchange(t *testing.T) {
	provider := crypto.NewStaticDataProvider()
	cryptoUsecase := usecase.NewCryptoUsecase(provider)

	tests := []struct {
		name     string
		from     string
		to       string
		amount   float64
		expected float64
		err      bool
	}{
		{
			name:     "valid USDT to BEER",
			from:     "USDT",
			to:       "BEER",
			amount:   1.0,
			expected: 40593.25477448192,
		},
		{
			name:     "valid WBTC to USDT",
			from:     "WBTC",
			to:       "USDT",
			amount:   1.0,
			expected: 57094.314314,
		},
		{
			name:   "unsupported currency",
			from:   "DOGE",
			to:     "USDT",
			amount: 10,
			err:    true,
		},
		{
			name:   "invalid amount",
			from:   "USDT",
			to:     "BEER",
			amount: -1,
			err:    true,
		},
		{
			name:   "empty fields",
			from:   "",
			to:     "USDT",
			amount: 1,
			err:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cryptoUsecase.Exchange(tt.from, tt.to, tt.amount)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.from, result.From)
				assert.Equal(t, tt.to, result.To)
				assert.InEpsilon(t, tt.expected, result.Amount, 1e-9)
			}
		})
	}
}
