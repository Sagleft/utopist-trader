package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) calcTPOrder(baseOrderPrice float64) order {
	strategySign := float64(-1)
	orderType := orderTypeBuy
	if b.isStrategyBuy() {
		orderType = orderTypeSell
		strategySign = 1
	}

	return order{
		Type:       orderType,
		PairSymbol: b.Config.PairSymbol,
		Qty:        b.Lap.CoinsQty,
		Price:      baseOrderPrice * (1 + strategySign*b.Config.ProfitPercent/100),
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
	// check order qty
	if o.Qty == 0 {
		return false, errors.New("order qty is 0")
	}

	// check order min deposit
	orderDeposit := o.Price * o.Qty
	if orderDeposit < b.PairMinDeposit {
		debug(
			"skip. the order deposit (%v) is not enough for the minimum: %v\n",
			orderDeposit, b.PairMinDeposit,
		)
		return false, nil
	}

	if o.Qty < b.PairData.MinAmount {
		debug(
			"order qty is too small (%v). minimum is %v. skip",
			o.Qty, b.PairData.MinAmount,
		)
		return false, nil
	}

	return true, nil
}

func (b *bot) checkTPOrder(o order) (bool, error) {
	return b.checkOrder(o)
}

// returns false when order doesn't fit
func (b *bot) checkMarketOrder(o order) (bool, error) {
	isOK, err := b.checkOrder(o)
	if err != nil {
		return false, err
	}
	if !isOK {
		return false, nil
	}

	// check available balance
	bl, err := b.getDepositBalance()
	if err != nil {
		return false, fmt.Errorf("get deposit balance: %w", err)
	}
	if bl.Balance < o.Price*o.Qty {
		log.Println("available deposit is not enought for the minimum order. skip")
		return false, nil
	}
	return true, nil
}

func (b *bot) calcMarketOrder() (order, error) {
	orderType := orderTypeSell
	if b.isStrategyBuy() {
		orderType = orderTypeBuy
	}

	priceRaw, err := b.getPairPrice()
	if err != nil {
		return order{}, err
	}

	debug("%s rate: %v\n", b.Config.PairSymbol, priceRaw)

	price := roundFloatFloor(priceRaw, b.PairData.RoundDealPrice)
	deposit := b.getOrderDeposit(price)
	qty := roundFloatFloor(deposit/price, b.PairData.RoundDealAmount)

	return order{
		Type:       orderType,
		PairSymbol: b.Config.PairSymbol,
		Qty:        qty,
		Price:      price,
	}, nil
}

func (b *bot) sendOrder(o order) (uexchange.OrderData, error) {
	var orderID int64
	var err error

	switch o.Type {
	default:
		return uexchange.OrderData{}, fmt.Errorf("unknown order type: %q", o.Type)
	case orderTypeBuy:
		orderID, err = b.Client.Buy(o.PairSymbol, o.Qty, o.Price)
	case orderTypeSell:
		orderID, err = b.Client.Sell(o.PairSymbol, o.Qty, o.Price)
	}

	if err != nil {
		return uexchange.OrderData{}, fmt.Errorf("send order %s: %w", o.ToString(), err)
	}

	// get placed order data
	return b.getOrderData(orderID)
}

func (b *bot) getOrderData(orderID int64) (uexchange.OrderData, error) {
	response, err := b.Client.GetOrderHistory(strconv.FormatInt(orderID, 10))
	if err != nil {
		return uexchange.OrderData{}, err
	}

	orderData := response.Order
	orderData.Amount = orderData.ExecutedValue / orderData.ExecutedPrice

	return orderData, nil
}

func (b *bot) getOrderDeposit(price float64) float64 {
	depositPercent := b.getIntervalDepositPercent(price)

	debug("use deposit percent: %v\n", depositPercent)

	return depositPercent * b.Config.Deposit / 100
}

func (b *bot) isTPPlaced() bool {
	return b.Lap.TPOrderID != 0
}
