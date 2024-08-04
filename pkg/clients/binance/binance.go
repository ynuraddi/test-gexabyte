package binance

import (
	"context"

	binance_connector "github.com/binance/binance-connector-go"
)

type Client interface {
	KlineService(ctx context.Context, symbol, interval string, startTime, endTime int64, limit int) ([]*binance_connector.KlinesResponse, error)
	TickerPriceService(ctx context.Context, symbol string) (*binance_connector.TickerPriceResponse, error)
	Ticker24hService(ctx context.Context, symbol string) (*binance_connector.Ticker24hrResponse, error)
}

type client struct {
	binance *binance_connector.Client
}

func New(cfg *Config) Client {
	c := binance_connector.NewClient(cfg.ApiKey, cfg.SecretKey)

	return &client{c}
}

func (c *client) KlineService(ctx context.Context, symbol, interval string, startTime, endTime int64, limit int) ([]*binance_connector.KlinesResponse, error) {
	res, err := c.binance.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(uint64(startTime)).
		EndTime(uint64(endTime)).
		Limit(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: add error handler

	return res, err
}

func (c *client) TickerPriceService(ctx context.Context, symbol string) (*binance_connector.TickerPriceResponse, error) {
	res, err := c.binance.NewTickerPriceService().Symbol(symbol).Do(ctx)

	return res, err
}
func (c *client) Ticker24hService(ctx context.Context, symbol string) (*binance_connector.Ticker24hrResponse, error) {
	res, err := c.binance.NewTicker24hrService().Symbol(symbol).Do(ctx)

	return res, err
}
