package service

import (
	"context"
	"gexabyte/internal/config"
	"gexabyte/internal/model"
	"gexabyte/internal/repository"
	"gexabyte/internal/service/currency"
	"gexabyte/pkg/clients/binance"
	"log/slog"
)

type Manager struct {
	Currency Currency
}

type Currency interface {
	// Symbol
	Create(ctx context.Context, symbol string) error
	List(ctx context.Context) ([]model.Currency, error)

	// Price
	CreatePrice(ctx context.Context, rates ...model.CurrencyPrice) error
	ListPrices(ctx context.Context) ([]model.CurrencyPrice, error)
	GetCurrentPrices(ctx context.Context, symbols ...string) ([]model.GetCurrencyPriceDTO, error)
	GetStat24H(ctx context.Context, symbols ...string) ([]model.GetCurrencyStat24HDTO, error)
	GetPriceHistorical(ctx context.Context, req model.GetCurrencyPriceHistoricalDTOReq) (*model.GetCurrencyPriceHistoricalDTORes, error)
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	binanceClient binance.Client,
	repository *repository.Manager,
) *Manager {
	currency := currency.NewCurrency(repository.Currency, repository.CurrencyPrice, binanceClient, logger)

	return &Manager{
		Currency: currency,
	}
}
