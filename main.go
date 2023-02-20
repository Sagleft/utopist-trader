package main

import (
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uexchange-go"
	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

func main() {
	swissknife.PrintIntroMessage(previewTitle, donateAddress)

	b := newBot()

	if err := swissknife.CheckErrors(
		b.parseConfig,
		b.verifyConfig,
		b.auth,
		b.verifyTradePair,
		b.loadPairData,
		b.verifyPairData,
		b.verifyBalance,
		b.runCheckExchangeCron,
	); err != nil {
		color.Red(err.Error())
		return
	}

	log.Println("bot started")
	swissknife.RunInBackground()
}

func newBot() *bot {
	return &bot{
		Client: uexchange.NewClient(),
	}
}

func (b *bot) auth() error {
	debug("connect to exchange..")

	if _, err := b.Client.Auth(uexchange.Credentials{
		AccountPublicKey: b.Config.Exchange.Pubkey,
		Password:         b.Config.Exchange.Password,
	}); err != nil {
		return err
	}

	success("exchange connected")
	return nil
}

func (b *bot) runCheckExchangeCron() error {
	debug("setup cron..")

	go simplecron.NewCronHandler(
		func() {
			if err := b.checkExchange(); err != nil {
				color.Red("check exchange: %s", err.Error())
			}
		},
		time.Duration(b.Config.IntervalTimeoutSeconds)*time.Second,
	).Run(b.Config.NoWait)

	success("crone initiated")
	return nil
}
