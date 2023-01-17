package main

import (
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
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
		b.run,
	); err != nil {
		log.Fatalln(err)
	}

	swissknife.RunInBackground()
}

func newBot() *bot {
	return &bot{}
}

func (b *bot) parseConfig() error {
	return swissknife.ParseStructFromJSONFile(configPath, &b.Config)
}

func (b *bot) run() error {
	simplecron.NewCronHandler(
		b.checkExchange,
		time.Duration(b.Config.CheckExchangeTimeout)*time.Second,
	).Run(checkExchangeAtStart)
	return nil
}

func (b *bot) checkExchange() {
	fmt.Println("timer tick")
}
