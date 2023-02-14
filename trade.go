package main

import (
	"github.com/Sagleft/uexchange-go"
)

func (b *bot) resetLap() {
	b.Lap = lap{}
}

func (b *bot) getIntervalDepositPercent(currentPrice float64) float64 {
	if b.isFirstInterval() {
		return b.Config.IntervalDepositMaxPercent
	}

	intervalMaxDeposit := b.Config.Deposit * b.Config.IntervalDepositMaxPercent / 100
	minPrice := b.Lap.LastPriceLevel * (1 - b.Config.IntervalDepositMaxPercent/2)
	maxPrice := b.Lap.LastPriceLevel * (1 + b.Config.IntervalDepositMaxPercent/2)

	if currentPrice < minPrice {
		return 100
	}
	if currentPrice > maxPrice {
		return 0
	}

	return 100 - (currentPrice-minPrice)/intervalMaxDeposit
}

func (b *bot) updateDepositUsed(orderData uexchange.OrderData) {
	b.Lap.CoinsQty += orderData.Amount
	b.Lap.DepositSpent += orderData.Value
}

func (b *bot) handleInterval() error {

	// calc market order
	defOrder, err := b.calcMarketOrder()
	if err != nil {
		return err
	}
	orderIsOK, err := b.checkOrder(defOrder)
	if err != nil {
		return err
	}
	if !orderIsOK {
		return nil
	}

	if b.isTPPlaced() {
		// cancel old TP order
		if err := b.Client.Cancel(b.Lap.TPOrderID); err != nil {
			return err
		}
	}

	// market buy
	orderData, err := b.sendOrder(defOrder)
	if err != nil {
		return err
	}

	b.updateDepositUsed(orderData)

	// place TP
	tpOrderID, err := b.sendTPOrder(orderData.Price)
	if err != nil {
		return err
	}

	// update TP order ID
	b.Lap.TPOrderID = tpOrderID
	return nil
}

func (b *bot) checkExchange() error {
	if b.Lap.IntervalNumber == 0 {
		return b.handleInterval()
	}

	// TODO: check TP
	// if executed -> lap finished
	// else:
	//    cancel TP,
	return b.handleInterval()
}
