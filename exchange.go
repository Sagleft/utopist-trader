package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Sagleft/uexchange-go"
)

func (b *bot) verifyTradePair() error {
	log.Println("verify trade pair..")

	if _, err := b.getTradePair(b.Config.PairSymbol); err != nil {
		return err
	}

	success("done")
	return nil
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

func (b *bot) getBalancePerTicker() (botPairBalance, error) {
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

func (b *bot) getDepositBalance() (botTickerBalance, error) {
	pairBalance, err := b.getBalancePerTicker()
	if err != nil {
		return botTickerBalance{}, err
	}

	switch b.Config.Strategy {
	default:
		return botTickerBalance{}, fmt.Errorf("unknown bot next action: %v", b.Config.Strategy)
	case botStrategyBuy:
		return pairBalance.QuoteAsset, nil
	case botStrategySell:
		return pairBalance.BaseAsset, nil
	}
}

func (b *bot) verifyBalance() error {
	log.Println("verify balance..")

	t, err := b.getDepositBalance()
	if err != nil {
		return err
	}

	if t.Balance < b.Config.Deposit {
		return fmt.Errorf(
			"%s balance not enough. available %v, needed %v",
			t.Ticker, t.Balance, b.Config.Deposit,
		)
	}
	success("done")
	return nil
}

// get pair price for market order
func (b *bot) getPairPrice() (float64, error) {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.PairSymbol))
	if err != nil {
		return 0, err
	}

	if b.isStrategyBuy() {
		return pairData.BestBidPrice, nil
	}
	return pairData.BestAskPrice, nil
}

func (b *bot) loadPairData() error {
	log.Println("load pair data..")

	pairs, err := b.Client.GetPairs()
	if err != nil {
		return err
	}

	for _, p := range pairs {
		if p.Pair.PairCode == b.Config.PairSymbol {
			b.PairData = p.Pair
			b.PairMinDeposit = p.Pair.MinPrice * p.Pair.MinAmount

			log.Printf("pair min deposit: %v\nmin amount: %v\n", b.PairMinDeposit, p.Pair.MinAmount)

			success("done")
			return nil
		}
	}
	return fmt.Errorf("pair %q not found", b.Config.PairSymbol)
}

func (b *bot) verifyPairData() error {
	log.Println("verify pair data..")

	if b.PairData.RoundDealAmount == 0 {
		return errors.New("load pair data: round deal amount is 0. data not available")
	}
	if b.PairData.RoundDealPrice == 0 {
		return errors.New("load pair data: round deal price is 0. data not available")
	}

	success("done")
	return nil
}
