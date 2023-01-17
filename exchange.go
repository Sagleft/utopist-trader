package main

import (
	"fmt"
	"strings"

	"github.com/Sagleft/uexchange-go"
	"github.com/fatih/color"
)

func (b *bot) verifyTradePair() error {
	_, err := b.getTradePair(b.Config.TradePair)
	return err
}

func (b *bot) getTradePair(code string) (uexchange.PairsDataContainer, error) {
	pairs, err := b.Client.GetPairs()
	if err != nil {
		return uexchange.PairsDataContainer{}, fmt.Errorf("get pairs: %w", err)
	}

	for _, p := range pairs {
		if strings.EqualFold(code, p.Pair.PairCode) {
			return p, nil
		}
	}
	return uexchange.PairsDataContainer{}, fmt.Errorf("%q trade pair not found", b.Config.TradePair)
}

func (b *bot) getPairParts() botPairData {
	pairParts := strings.Split(strings.ToLower(b.Config.TradePair), "_")
	return botPairData{
		BaseAsset:  pairParts[0],
		QuoteAsset: pairParts[1],
	}
}

func (b *bot) getBalance() (botPairBalance, error) {
	balances, err := b.Client.GetBalance()
	if err != nil {
		return botPairBalance{}, fmt.Errorf("get balance: %w", err)
	}

	pairParts := b.getPairParts()

	r := botPairBalance{}
	for _, balanceData := range balances {
		if balanceData.Currency.Name == pairParts.BaseAsset {
			r.BaseAsset = botTickerBalance{
				Ticker:  pairParts.BaseAsset,
				Balance: balanceData.Balance,
			}
		}
		if balanceData.Currency.Name == pairParts.QuoteAsset {
			r.QuoteAsset = botTickerBalance{
				Ticker:  pairParts.QuoteAsset,
				Balance: balanceData.Balance,
			}
		}
	}

	return r, nil
}

func (b *bot) checkBalance() error {
	pairBalance, err := b.getBalance()
	if err != nil {
		return err
	}

	var t botTickerBalance
	if b.Config.StartFromBuy {
		t = pairBalance.QuoteAsset
	} else {
		t = pairBalance.BaseAsset
	}

	if t.Balance < b.Config.Deposit {
		return fmt.Errorf(
			"%s balance not enough. available %v, needed %v",
			t.Ticker, t.Balance, b.Config.Deposit,
		)
	}

	return nil
}

func (b *bot) checkExchange() {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.TradePair))
	if err != nil {
		color.Red(err.Error())
		return
	}

	fmt.Printf("ask %v, bid %v\n", pairData.BestAskPrice, pairData.BestBidPrice)
}
