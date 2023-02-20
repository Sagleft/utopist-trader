package main

import (
	"sync"

	"github.com/Sagleft/uexchange-go"
)

type bot struct {
	Client         *uexchange.Client
	Config         config
	PairData       uexchange.PairData
	PairMinDeposit float64
	Lap            lap

	HandleIntervalLock sync.Mutex
}

type lap struct {
	Number         int
	IntervalNumber int
	LastPriceLevel float64 // previous market order price
	TPOrderID      int64

	CoinsQty                   float64
	DepositSpent               float64
	TPAlreadyPartiallyExecuted bool
}

type config struct {
	Strategy                  string         `json:"strategy"`
	IntervalDepositMaxPercent float64        `json:"intervalDepositMaxPercent"`
	ProfitPercent             float64        `json:"profitPercent"`
	PairSymbol                string         `json:"tradePair"`
	Deposit                   float64        `json:"deposit"`
	Exchange                  exchangeConfig `json:"exchange"`
	IntervalTimeoutSeconds    float64        `json:"intervalTimeoutSeconds"`
	NoWait                    bool           `json:"noWait"`
	IsDebug                   bool           `json:"debug"`
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

type order struct {
	Type       string  `json:"type"`
	PairSymbol string  `json:"symbol"`
	Qty        float64 `json:"amount"`
	Price      float64 `json:"price"`
}

type orderExecutedState struct {
	IsFullExecuted      bool
	IsPartiallyExecuted bool
	Data                uexchange.OrderData
}
