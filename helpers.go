package main

import (
	"encoding/json"
	"errors"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

func (b *bot) parseConfig() error {
	log.Println("parse config..")
	if err := swissknife.ParseStructFromJSONFile(configPath, &b.Config); err != nil {
		return err
	}

	success("done")
	return nil
}

func (b *bot) verifyConfig() error {
	log.Println("verify config..")

	if b.Config.IntervalDepositMaxPercent == 0 {
		return errors.New("invalid `intervalDepositMaxPercent`: value must be set")
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

func (b *bot) setLastPriceLevel(price float64) {
	b.Lap.LastPriceLevel = price
}

func success(info string, a ...interface{}) {
	if info == "" {
		return
	}
	color.Green(" [ "+info+" ]", a...)
}

func roundFloatFloor(val float64, precision int) float64 {
	f, _ := decimal.NewFromFloat(val).RoundFloor(int32(precision)).Float64()
	return f
}
