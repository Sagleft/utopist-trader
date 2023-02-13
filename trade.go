package main

import (
	"fmt"
	"strings"
)

func (b *bot) startNewLap() error {
	if err := b.checkBalance(); err != nil {
		return err
	}

	if err := b.checkExchange(); err != nil {
		return err
	}

	// TODO
	return nil
}

func (b *bot) checkExchange() error {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.TradePair))
	if err != nil {
		return err
	}

	fmt.Printf("ask %v, bid %v\n", pairData.BestAskPrice, pairData.BestBidPrice)

	// TODO
	return nil
}
