package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func (b *bot) checkExchange() {
	pairData, err := b.Client.GetPairPrice(strings.ToLower(b.Config.TradePair))
	if err != nil {
		color.Red(err.Error())
		return
	}

	fmt.Printf("ask %v, bid %v\n", pairData.BestAskPrice, pairData.BestBidPrice)
}
