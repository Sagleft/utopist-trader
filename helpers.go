package main

import (
	"encoding/json"

	swissknife "github.com/Sagleft/swiss-knife"
)

func (b *bot) parseConfig() error {
	return swissknife.ParseStructFromJSONFile(configPath, &b.Config)
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
