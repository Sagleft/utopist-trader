package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) sendTPOrder(baseOrderPrice float64) (int64, error) {
	tpOrderData, err := b.sendOrder(order{
		PairSymbol: b.Config.PairSymbol,
		Qty:        b.Lap.CoinsQty,
		Price:      baseOrderPrice * (1 + b.Config.ProfitPercent/100),
	})
	if err != nil {
		return 0, err
	}

	return tpOrderData.OrderID, nil
}

func (b *bot) sendMarketOrder() (uexchange.OrderData, error) {
	o, err := b.calcMarketOrder()
	if err != nil {
		return uexchange.OrderData{}, err
	}

	// check order min deposit
	orderDeposit := o.Price * o.Qty
	if orderDeposit < b.PairMinDeposit {
		log.Printf(
			"skip. the order deposit (%v) is not enough for the minimum: %v\n",
			orderDeposit, b.PairMinDeposit,
		)
		return uexchange.OrderData{}, nil
	}

	// check available balance
	bl, err := b.getDepositBalance()
	if err != nil {
		return uexchange.OrderData{}, err
	}
	if bl.Balance < orderDeposit {
		log.Println("available deposit is not enought for the minimum order. skip")
		return uexchange.OrderData{}, nil
	}

	return b.sendOrder(o)
}

func (b *bot) calcMarketOrder() (order, error) {
	price, err := b.getPairPrice()
	if err != nil {
		return order{}, err
	}

	deposit, err := b.getOrderDeposit(price)
	if err != nil {
		return order{}, err
	}

	return order{
		PairSymbol: b.Config.PairSymbol,
		Qty:        deposit / price,
		Price:      price,
	}, nil
}

func isOrderEmpty(orderData uexchange.OrderData) bool {
	return orderData.OrderID == 0
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
	orderData, err := b.Client.GetOrderHistory(strconv.FormatInt(orderID, 10))
	if err != nil {
		return uexchange.OrderData{}, err
	}
	return orderData.Order, nil
}

func (b *bot) getOrderDeposit(price float64) (float64, error) {
	return b.getIntervalDepositPercent(price), nil
}
