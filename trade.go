package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (b *bot) resetLap() {
	b.Lap = lap{}
}

// get pair price for market order
func (b *bot) getPairPrice() (float64, error) {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.TradePair))
	if err != nil {
		return 0, err
	}

	if b.isStrategyBuy() {
		return pairData.BestAskPrice, nil
	}
	return pairData.BestBidPrice, nil
}

func (b *bot) getOrderDeposit() (float64, error) {
	price, err := b.getPairPrice()
	if err != nil {
		return 0, err
	}

	return b.getIntervalDepositPercent(price), nil
}

func (b *bot) getIntervalDepositPercent(currentPrice float64) float64 {
	intervalMaxDeposit := b.Config.Deposit * b.Config.IntervalDepositMaxPercent / 100
	minPrice := currentPrice * (1 - b.Config.IntervalDepositMaxPercent/2)
	maxPrice := currentPrice * (1 + b.Config.IntervalDepositMaxPercent/2)

	if currentPrice < minPrice {
		return 100
	}
	if currentPrice > maxPrice {
		return 0
	}

	return 100 - (currentPrice-minPrice)/intervalMaxDeposit
}

type order struct {
	PairSymbol string
	Qty        float64
	Price      float64
}

func (o order) ToString() string {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func (b *bot) isStrategyBuy() bool {
	return b.Config.Strategy == botStrategyBuy
}

func (b *bot) sendOrder(o order) error {
	var err error
	if b.isStrategyBuy() {
		_, err = b.Client.Buy(o.PairSymbol, o.Qty, o.Price)
	} else {
		_, err = b.Client.Sell(o.PairSymbol, o.Qty, o.Price)
	}
	if err != nil {
		return fmt.Errorf("send market order %s: %w", o.ToString(), err)
	}
	return nil
}

func (b *bot) checkExchange() error {
	if b.Lap.IntervalNumber == 0 {
		// market buy
		if err := b.sendOrder(order{}); err != nil {
			return err
		}

		// TODO: place TP
		return nil
	}

	// TODO: check TP
	// if executed -> lap finished
	// else:
	//    cancel TP,
	//    market buy
	if err := b.sendOrder(order{}); err != nil {
		return err
	}
	return nil
}
