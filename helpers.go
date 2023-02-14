package main

import (
	"encoding/json"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/fatih/color"
)

func (b *bot) parseConfig() error {
	log.Println("parse config..")
	if err := swissknife.ParseStructFromJSONFile(configPath, &b.Config); err != nil {
		return err
	}

	success("done")
	return nil
}

func (o order) ToString() string {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func (b *bot) isStrategyBuy() bool {
	return b.Config.Strategy == botStrategyBuy
}

func (b *bot) isFirstInterval() bool {
	return b.Lap.IntervalNumber == 0
}

func success(info string, a ...interface{}) {
	if info == "" {
		return
	}
	color.Green("[ "+info+" ]", a...)
}
