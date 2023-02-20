package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIntervalDepositPercent(t *testing.T) {
	b := bot{
		Config: config{
			IntervalDepositMaxPercent: 5,
			Deposit:                   100,
		},
		Lap: lap{
			LastPriceLevel: 0.6,
		},
	}

	depositPercent := b.getIntervalDepositPercent(0.6252)

	assert.NotEqual(t, 0, depositPercent)
}
