package main

import "github.com/Sagleft/uexchange-go"

type bot struct {
	Client *uexchange.Client
	Config config
}

type config struct {
	DipTresholdPercent   float64        `json:"dipTresholdPercent"`
	UpwardTrend          percentParam   `json:"upwardTrend"`
	ProfitPercent        float64        `json:"profitPercent"`
	StopLoss             percentParam   `json:"stopLoss"`
	TradePair            string         `json:"tradePair"`
	Deposit              float64        `json:"deposit"`
	Exchange             exchangeConfig `json:"exchange"`
	CheckExchangeTimeout float64        `json:"checkExchangeTimeoutSeconds"`
	StartFromBuy         bool           `json:"startFromBuyAction"`
}

type percentParam struct {
	Value   float64 `json:"percent"`
	Enabled bool    `json:"enabled"`
}

type exchangeConfig struct {
	Pubkey   string `json:"pubkey"`
	Password string `json:"password"`
}
