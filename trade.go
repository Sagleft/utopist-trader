package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) resetLap() {
	b.Lap = lap{}
}

func (b *bot) loadPairData() error {
	pairs, err := b.Client.GetPairs()
	if err != nil {
		return err
	}

	for _, p := range pairs {
		if p.Pair.PairCode == b.Config.PairSymbol {
			b.PairData = p.Pair
			b.PairMinDeposit = p.Pair.MinPrice * p.Pair.MinAmount
			return nil
		}
	}
	return fmt.Errorf("pair %q not found", b.Config.PairSymbol)
}

// get pair price for market order
func (b *bot) getPairPrice() (float64, error) {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.PairSymbol))
	if err != nil {
		return 0, err
	}

	if b.isStrategyBuy() {
		return pairData.BestAskPrice, nil
	}
	return pairData.BestBidPrice, nil
}

func (b *bot) getOrderDeposit(price float64) (float64, error) {
	return b.getIntervalDepositPercent(price), nil
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

func (b *bot) updateDepositUsed(orderData uexchange.OrderData) {
	b.Lap.CoinsQty += orderData.Amount
	b.Lap.DepositSpent += orderData.Value
}

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

func (b *bot) handleInterval() error {
	// market buy
	orderData, err := b.sendMarketOrder()
	if err != nil {
		return err
	}
	if isOrderEmpty(orderData) {
		return nil // skip
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
