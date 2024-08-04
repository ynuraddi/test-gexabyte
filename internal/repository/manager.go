package repository

import (
	"context"

	"gexabyte/internal/config"
	"gexabyte/internal/model"
	repo "gexabyte/internal/repository/postgres"
	"gexabyte/pkg/clients/postgres"
)

type Manager struct {
	Currency      Currency
	CurrencyPrice CurrencyPrice
}

type Currency interface {
	Create(ctx context.Context, symbol string) error
	GetBySymbol(ctx context.Context, symbol string) (model.Currency, error)
	List(ctx context.Context) ([]model.Currency, error)
}

type CurrencyPrice interface {
	Create(ctx context.Context, rates ...model.CurrencyPrice) error
	List(ctx context.Context) ([]model.CurrencyPrice, error)
}

func NewRepository(cfg *config.Config) (*Manager, error) {
	dbClient, err := postgres.NewClient(postgres.Config{DSN: cfg.Postgres.DSN})
	if err != nil {
		return nil, err
	}

	if err := postgres.RunDBMigration(cfg.Postgres.MigrationURL, cfg.Postgres.DSN); err != nil {
		return nil, err
	}

	currencyPairs := repo.NewCurrency(dbClient.DB)
	currencyPrice := repo.NewCurrencyPrice(dbClient.DB)

	return &Manager{
		Currency:      currencyPairs,
		CurrencyPrice: currencyPrice,
	}, nil
}
