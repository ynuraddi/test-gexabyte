package currency

import (
	"context"
	"gexabyte/internal/model"
	"gexabyte/internal/repository"
	"gexabyte/pkg/clients/binance"
	"log/slog"
	"time"
)

const LoggerGroup = "CurrencyService"

type Currency struct {
	currencyRepo      repository.Currency
	currencyPriceRepo repository.CurrencyPrice

	binanceClient binance.Client

	logger *slog.Logger

	priceCheckTicker   *time.Ticker
	priceCheckInterval time.Duration
}

func NewCurrency(
	currencyRepo repository.Currency,
	currencyPriceRepo repository.CurrencyPrice,
	binanceClient binance.Client,
	logger *slog.Logger,
) *Currency {
	return &Currency{
		currencyRepo:      currencyRepo,
		currencyPriceRepo: currencyPriceRepo,

		binanceClient: binanceClient,

		logger: logger.WithGroup(LoggerGroup),

		priceCheckTicker:   time.NewTicker(10 * time.Minute),
		priceCheckInterval: 10 * time.Minute,
	}
}

func (s *Currency) Create(ctx context.Context, symbol string) error {
	return s.currencyRepo.Create(ctx, symbol)
}

func (s *Currency) List(ctx context.Context) ([]model.Currency, error) {
	return s.currencyRepo.List(ctx)
}
