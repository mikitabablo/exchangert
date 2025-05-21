package domain

type CryptoCurrency struct {
	Name          string
	RateUSD       float64
	DecimalPlaces int
}

type CryptoExchangeResult struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
