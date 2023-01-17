package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uexchange-go"
	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

const (
	configPath           = "config.json"
	checkExchangeAtStart = true
)

func main() {
	b := newBot()

	if err := swissknife.CheckErrors(
		b.parseConfig,
		b.auth,
		b.verifyTradePair,
		b.run,
	); err != nil {
		log.Fatalln(err)
	}

	swissknife.RunInBackground()
}

func newBot() *bot {
	return &bot{
		Client: uexchange.NewClient(),
	}
}

func (b *bot) parseConfig() error {
	return swissknife.ParseStructFromJSONFile(configPath, &b.Config)
}

func (b *bot) auth() error {
	_, err := b.Client.Auth(uexchange.Credentials{
		AccountPublicKey: b.Config.Exchange.Pubkey,
		Password:         b.Config.Exchange.Password,
	})
	return err
}

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

func (b *bot) run() error {
	simplecron.NewCronHandler(
		b.checkExchange,
		time.Duration(b.Config.CheckExchangeTimeout)*time.Second,
	).Run(checkExchangeAtStart)
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
