package model

type Currency struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

type CurrencyPrice struct {
	ID         int
	CurrencyID int
	Price      float64
	Time       int64
}
