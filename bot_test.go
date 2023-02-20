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
			IntervalNumber: 1,
			LastPriceLevel: 0.6,
		},
	}

	currentPrice := float64(0.6252)
	depositPercent := b.getIntervalDepositPercent(currentPrice)

	assert.Equal(t, 3.4748, depositPercent)
}
