package main

import "github.com/Sagleft/uexchange-go"

type bot struct {
	Client *uexchange.Client
	Config config
	Lap    lap
}

type lap struct {
	Number         int
	IntervalNumber int
	LastPriceLevel float64
	Position       position
}

type config struct {
	Strategy                  string         `json:"strategy"`
	IntervalDepositMaxPercent float64        `json:"intervalDepositMaxPercent"`
	ProfitPercent             float64        `json:"profitPercent"`
	TradePair                 string         `json:"tradePair"`
	Deposit                   float64        `json:"deposit"`
	Exchange                  exchangeConfig `json:"exchange"`
	IntervalTimeoutSeconds    float64        `json:"intervalTimeoutSeconds"`
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
