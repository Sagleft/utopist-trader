package main

import (
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
)

const (
	configPath = "config.json"
)

func main() {
	b := newBot()

	if err := swissknife.CheckErrors(
		b.parseConfig,
		b.run,
	); err != nil {
		log.Fatalln(err)
	}
}

func newBot() *bot {
	return &bot{}
}

func (b *bot) parseConfig() error {
	return swissknife.ParseStructFromJSONFile(configPath, &b.Config)
}

func (b *bot) run() error {
	// TODO
	return nil
}
