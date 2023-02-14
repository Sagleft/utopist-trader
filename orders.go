package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) calcTPOrder(baseOrderPrice float64) order {
	return order{
		PairSymbol: b.Config.PairSymbol,
		Qty:        b.Lap.CoinsQty,
		Price:      baseOrderPrice * (1 + b.Config.ProfitPercent/100),
	}
}

func (b *bot) sendTPOrder(o order) (int64, error) {
	tpOrderData, err := b.sendOrder(o)
	if err != nil {
		return 0, err
	}

	return tpOrderData.OrderID, nil
}

// returns false when order doesn't fit
func (b *bot) checkOrder(o order) (bool, error) {
	// check order min deposit
	orderDeposit := o.Price * o.Qty
	if orderDeposit < b.PairMinDeposit {
		log.Printf(
			"skip. the order deposit (%v) is not enough for the minimum: %v\n",
			orderDeposit, b.PairMinDeposit,
		)
		return false, nil
	}

	// check available balance
	bl, err := b.getDepositBalance()
	if err != nil {
		return false, err
	}
	if bl.Balance < orderDeposit {
		log.Println("available deposit is not enought for the minimum order. skip")
		return false, nil
	}
	return true, nil
}

func (b *bot) calcMarketOrder() (order, error) {
	price, err := b.getPairPrice()
	if err != nil {
		return order{}, err
	}

	deposit := b.getOrderDeposit(price)

	return order{
		PairSymbol: b.Config.PairSymbol,
		Qty:        deposit / price,
		Price:      price,
	}, nil
}

func (b *bot) sendOrder(o order) (uexchange.OrderData, error) {
	var orderID int64
	var err error
	if b.isStrategyBuy() {
		orderID, err = b.Client.Buy(o.PairSymbol, o.Qty, o.Price)
	} else {
		orderID, err = b.Client.Sell(o.PairSymbol, o.Qty, o.Price)
	}
	if err != nil {
		return uexchange.OrderData{}, fmt.Errorf("send market order %s: %w", o.ToString(), err)
	}

	// get placed order data
	return b.getOrderData(orderID)
}

func (b *bot) getOrderData(orderID int64) (uexchange.OrderData, error) {
	orderData, err := b.Client.GetOrderHistory(strconv.FormatInt(orderID, 10))
	if err != nil {
		return uexchange.OrderData{}, err
	}
	return orderData.Order, nil
}

func (b *bot) getOrderDeposit(price float64) float64 {
	return b.getIntervalDepositPercent(price) * b.Config.Deposit
}

func (b *bot) isTPPlaced() bool {
	return b.Lap.TPOrderID == 0
}
