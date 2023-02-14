package main

import (
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uexchange-go"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

func main() {
	figure.NewColorFigure(" utopist-trader $$$", "", "green", true).
		Scroll(3*1000, 200, "left")

	b := newBot()

	if err := swissknife.CheckErrors(
		b.parseConfig,
		b.auth,
		b.verifyTradePair,
		b.loadPairData,
		b.verifyBalance,
		b.runCheckExchangeCron,
	); err != nil {
		color.Red(err.Error())
		return
	}

	swissknife.RunInBackground()
}

func newBot() *bot {
	return &bot{
		Client: uexchange.NewClient(),
	}
}

func (b *bot) auth() error {
	log.Println("connect to exchange..")

	if _, err := b.Client.Auth(uexchange.Credentials{
		AccountPublicKey: b.Config.Exchange.Pubkey,
		Password:         b.Config.Exchange.Password,
	}); err != nil {
		return err
	}

	success("connected")
	return nil
}

func (b *bot) runCheckExchangeCron() error {
	log.Println("setup cron..")

	go simplecron.NewCronHandler(
		func() {
			if err := b.checkExchange(); err != nil {
				color.Red("check exchange: %s", err.Error())
			}
		},
		time.Duration(b.Config.IntervalTimeoutSeconds)*time.Second,
	).Run(checkExchangeAtStart)

	success("done")
	return nil
}
