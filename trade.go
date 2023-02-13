package main

import "strings"

func (b *bot) resetLap() {
	b.Lap = lap{}
}

// get pair price for market order
func (b *bot) getPairPrice() (float64, error) {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.TradePair))
	if err != nil {
		return 0, err
	}

	if b.Config.Strategy == botStrategyBuy {
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

func (b *bot) checkExchange() error {
	if b.Lap.IntervalNumber == 0 {
		// TODO: market buy
		// TODO: place TP
		return nil
	}

	// TODO: check TP
	// if executed -> lap finished
	// else: cancel TP, market buy
	return nil
}
