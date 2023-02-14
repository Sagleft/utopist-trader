package main

import (
	"fmt"
	"strings"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) verifyTradePair() error {
	_, err := b.getTradePair(b.Config.PairSymbol)
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
	return uexchange.PairsDataContainer{}, fmt.Errorf("%q trade pair not found", b.Config.PairSymbol)
}

func (b *bot) getPairParts() botPairData {
	pairParts := strings.Split(strings.ToLower(b.Config.PairSymbol), "_")
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
	switch b.Config.Strategy {
	default:
		return fmt.Errorf("unknown bot next action: %v", b.Config.Strategy)
	case botStrategyBuy:
		t = pairBalance.QuoteAsset
	case botStrategySell:
		t = pairBalance.BaseAsset
	}

	if t.Balance < b.Config.Deposit {
		return fmt.Errorf(
			"%s balance not enough. available %v, needed %v",
			t.Ticker, t.Balance, b.Config.Deposit,
		)
	}

	b.Lap.Position.AvailableDeposit = b.Config.Deposit
	return nil
}
