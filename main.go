package main

import (
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uexchange-go"
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

func (b *bot) run() error {
	simplecron.NewCronHandler(
		b.checkExchange,
		time.Duration(b.Config.CheckExchangeTimeout)*time.Second,
	).Run(checkExchangeAtStart)
	return nil
}

func (b *bot) checkExchange() {

}
