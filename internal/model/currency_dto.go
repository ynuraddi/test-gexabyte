package model

type GetCurrencyPriceDTO struct {
	Symbol string
	Price  float64
	Time   int64
}

type GetCurrencyStat24HDTO struct {
	Symbol    string  `json:"symbol"`
	OpenPrice float64 `json:"open_price"`
	LastPrice float64 `json:"last_price"`
	HighPrice float64 `json:"high_price"`
	LowPrice  float64 `json:"low_price"`
	// Volume      float64 `json:"volume,string"`
	// QuoteVolume float64 `json:"quoteVolume,string"`
	OpenTime  int64 `json:"open_time"`
	CloseTime int64 `json:"close_time"`
	// FirstID   int64 `json:"firstId"`
	// LastID    int64 `json:"lastId"`
	// Count     int   `json:"count"`
}

type GetCurrencyPriceHistoricalDTOReq struct {
	Symbol    string
	Interval  string
	StartTime int64
	EndTime   int64
	Limit     int
	Page      int
}

type CurrencyPriceInterval struct {
	OpenPrice  float64 `json:"open_price"`
	ClosePrice float64 `json:"close_price"`
	HighPrice  float64 `json:"high_price"`
	LowPrice   float64 `json:"low_price"`

	OpenTime  int64 `json:"open_time"`
	CloseTime int64 `json:"close_time"`
}

type GetCurrencyPriceHistoricalDTORes struct {
	Page    int `json:"page"`
	MaxPage int `json:"max_page"`

	Prices []CurrencyPriceInterval `json:"prices"`
}
