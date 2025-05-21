package openexchangerates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mikitabablo/exchangert/internal/domain"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetRates(base string, symbols []string) (*domain.RatesResponse, error) {
	url := fmt.Sprintf(
		"%s/latest.json?app_id=%s&base=%s&symbols=%s",
		c.baseURL,
		c.apiKey,
		base,
		strings.Join(symbols, ","),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ratesResp domain.RatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ratesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ratesResp, nil
}
