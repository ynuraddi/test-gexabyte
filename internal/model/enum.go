package model

import (
	"time"
)

type klineInterval struct {
	intervalDuration map[string]time.Duration
}

// Please note that if you are going to work with a month, then use time.AddDate()
// ref: https://developers.binance.com/docs/binance-spot-api-docs/rest-api#klinecandlestick-data
var KlineInterval klineInterval = klineInterval{
	intervalDuration: map[string]time.Duration{
		"1s":  1 * time.Second,
		"1m":  1 * time.Minute,
		"3m":  3 * time.Minute,
		"5m":  5 * time.Minute,
		"15m": 15 * time.Minute,
		"30m": 30 * time.Minute,
		"1h":  1 * time.Hour,
		"2h":  2 * time.Hour,
		"4h":  4 * time.Hour,
		"6h":  6 * time.Hour,
		"12h": 12 * time.Hour,
		"1d":  24 * time.Hour,
		"3d":  3 * 24 * time.Hour,
		"1w":  7 * 24 * time.Hour,
		"1M":  30 * 24 * time.Hour,
	},
}

func (k klineInterval) IsCorrect(interval string) bool {
	_, ok := k.intervalDuration[interval]
	return ok
}

func (k klineInterval) GetDuration(interval string) time.Duration {
	return k.intervalDuration[interval]
}
