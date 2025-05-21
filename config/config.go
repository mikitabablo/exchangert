package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server            Server
	OpenExchangeRates OpenExchangeRates
}

type OpenExchangeRates struct {
	Url   string // `env:"OPEN_EXCHANGE_RATES_URL"`
	AppId string // `env:"OPEN_EXCHANGE_RATES_APP_ID"`
}

type Server struct {
	Port string
}

func LoadConfig() (*Config, error) {
//	err := godotenv.Load("../../.env")
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env: %v", err)
	}

	cfg := &Config{}

	cfg.Server.Port = os.Getenv("HTTP_PORT")
	if cfg.Server.Port == "" {
		return nil, errors.New("no port provided")
	}

	cfg.OpenExchangeRates.Url = os.Getenv("OPEN_EXCHANGE_RATES_URL")
	if cfg.OpenExchangeRates.Url == "" {
		return nil, errors.New("no exchange rates url")
	}
	cfg.OpenExchangeRates.AppId = os.Getenv("OPEN_EXCHANGE_RATES_APP_ID")
	if cfg.OpenExchangeRates.Url == "" {
		return nil, errors.New("no exchange rates app id")
	}

	return cfg, nil
}
