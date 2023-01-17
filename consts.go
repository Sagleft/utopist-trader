package main

const (
	configPath           = "config.json"
	checkExchangeAtStart = true
)

type ActionType int

const (
	ActionBUY ActionType = iota
	ActionSELL
)
