package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

func (b *bot) parseConfig() error {
	if err := swissknife.ParseStructFromJSONFile(configPath, &b.Config); err != nil {
		return err
	}

	debugMode = b.Config.IsDebug
	debug("parse config")
	success("config parsed")
	return nil
}

func (b *bot) verifyConfig() error {
	debug("verify config..")

	if b.Config.IntervalDepositMaxPercent == 0 {
		return errors.New("invalid `intervalDepositMaxPercent`: value must be set")
	}

	success("config verified")
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
	color.Green(" [ "+info+" ]\n", a...)
}

func roundFloatFloor(val float64, precision int) float64 {
	f, _ := decimal.NewFromFloat(val).RoundFloor(int32(precision)).Float64()
	return f
}

func warn(info string, a ...interface{}) {
	color.Yellow("[WARN] "+info, a...)
}

func toJSON(v any) string {
	dataBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to json convert: %s"}`, err.Error())
	}
	return string(dataBytes)
}

func debug(info string, a ...any) {
	if !debugMode {
		return
	}

	log.Printf(info+"\n", a...)
}
