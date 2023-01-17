package main

import (
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uexchange-go"
	"github.com/fatih/color"
	simplecron "github.com/sagleft/simple-cron"
)

func main() {
	b := newBot()

	if err := swissknife.CheckErrors(
		b.parseConfig,
		b.initTradeAction,
		b.auth,
		b.verifyTradePair,
		b.checkBalance,
		b.run,
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

func (b *bot) initTradeAction() error {
	if b.Config.StartFromBuy {
		b.NextAction = ActionBUY
	} else {
		b.NextAction = ActionSELL
	}
	return nil
}

func (b *bot) auth() error {
	_, err := b.Client.Auth(uexchange.Credentials{
		AccountPublicKey: b.Config.Exchange.Pubkey,
		Password:         b.Config.Exchange.Password,
	})
	return err
}

func (b *bot) run() error {
	simplecron.NewCronHandler(
		b.checkExchange,
		time.Duration(b.Config.CheckExchangeTimeout)*time.Second,
	).Run(checkExchangeAtStart)
	return nil
}
