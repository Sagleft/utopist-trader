package main

import (
	"log"

	"github.com/Sagleft/uexchange-go"
	"github.com/fatih/color"
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
		"handle interval %v..\n",
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
		log.Printf("cancel old TP order: %v\n", b.Lap.TPOrderID)
		if err := b.Client.Cancel(b.Lap.TPOrderID); err != nil {
			return err
		}
	}

	// market buy
	log.Printf("send order: %s\n", defOrder.ToString())
	orderData, err := b.sendOrder(defOrder)
	if err != nil {
		return err
	}
	success("order placed: %v", orderData.OrderID)

	b.updateDepositUsed(orderData)

	// place TP
	log.Println("calc TP order..")
	tpOrder := b.calcTPOrder(orderData.Price)

	log.Printf("place TP order: %s\n", tpOrder.ToString())
	tpOrderID, err := b.sendTPOrder(tpOrder)
	if err != nil {
		return err
	}
	success("TP placed: %v", tpOrderID)

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
		Data:                tpOrderData,
	}, nil
}

func (b *bot) getLapProfit(tpState orderExecutedState) (float64, error) {
	return tpState.Data.ExecutedValue - b.Config.Deposit, nil
}

func (b *bot) markTPPartiallyExecuted() {
	b.Lap.TPAlreadyPartiallyExecuted = true
}

func (b *bot) checkExchange() error {
	if !b.HandleIntervalLock.TryLock() {
		return nil // prevent parallel run
	}
	defer b.HandleIntervalLock.Unlock()

	if b.Lap.IntervalNumber == 0 {
		return b.handleInterval()
	}

	tpState, err := b.getTPExecutedState()
	if err != nil {
		return err
	}

	if tpState.IsPartiallyExecuted && !b.Lap.TPAlreadyPartiallyExecuted {
		b.markTPPartiallyExecuted()
		color.HiYellow("TP was partially executed")
	}

	if tpState.IsFullExecuted {
		lapProfit, err := b.getLapProfit(tpState)
		if err != nil {
			return err
		}

		success("💰 Lap finished! Profit: %v\n", lapProfit)
		b.resetLap()
	}

	return b.handleInterval()
}
