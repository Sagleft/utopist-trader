package main

import "github.com/Sagleft/uexchange-go"

type bot struct {
	Client *uexchange.Client
	Config config

	Position position
}

type config struct {
	Action                      string         `json:"action"`
	ProfitPercent               float64        `json:"profitPercent"`
	StopLoss                    percentParam   `json:"stopLoss"`
	TradePair                   string         `json:"tradePair"`
	Deposit                     float64        `json:"deposit"`
	Exchange                    exchangeConfig `json:"exchange"`
	CheckExchangeTimeoutSeconds float64        `json:"checkExchangeTimeoutSeconds"`
	StartFromBuy                bool           `json:"startFromBuyAction"`
}

type percentParam struct {
	Value   float64 `json:"percent"`
	Enabled bool    `json:"enabled"`
}

type exchangeConfig struct {
	Pubkey   string `json:"pubkey"`
	Password string `json:"password"`
}

type botPairData struct {
	BaseAsset  string
	QuoteAsset string
}

type botPairBalance struct {
	BaseAsset  botTickerBalance
	QuoteAsset botTickerBalance
}

type botTickerBalance struct {
	Ticker  string
	Balance float64
}

type position struct {
	DepositSpent     float64
	InitialDeposit   float64
	AvailableDeposit float64
}
