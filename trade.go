package main

import (
	"log"

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

func (b *bot) incrementIntervalNumber() {
	b.Lap.IntervalNumber++
}

func (b *bot) handleInterval() error {
	log.Printf(
		"handle interval #v..\n",
		b.Lap.IntervalNumber,
	)

	defer b.incrementIntervalNumber()

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

func (b *bot) getTPExecutedState() (orderExecutedState, error) {
	tpOrderData, err := b.getOrderData(b.Lap.TPOrderID)
	if err != nil {
		return orderExecutedState{}, err
	}

	return orderExecutedState{
		IsFullExecuted:      tpOrderData.ExecutedValue == tpOrderData.OriginalValue,
		IsPartiallyExecuted: tpOrderData.ExecutedValue > 0,
	}, nil
}

func (b *bot) getLapProfit() (float64, error) {
	// TODO

	return 0, nil
}

func (b *bot) checkExchange() error {
	if !b.HandleIntervalLock.TryLock() {
		return nil // prevent parallel run
	}

	if b.Lap.IntervalNumber == 0 {
		return b.handleInterval()
	}

	tpState, err := b.getTPExecutedState()
	if err != nil {
		return err
	}

	if tpState.IsFullExecuted {
		lapProfit, err := b.getLapProfit()
		if err != nil {
			return err
		}

		log.Printf("Wow! Lap finished! Profit: %v\n", lapProfit)
		b.resetLap()
	}

	return b.handleInterval()
}
