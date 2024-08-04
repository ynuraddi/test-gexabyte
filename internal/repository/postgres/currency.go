package postgres

import (
	"context"
	"database/sql"
	"gexabyte/internal/model"
)

type CurrencyRepo struct {
	db *sql.DB
}

func NewCurrency(db *sql.DB) *CurrencyRepo {
	return &CurrencyRepo{
		db: db,
	}
}

func (r *CurrencyRepo) Create(ctx context.Context, symbol string) error {
	query := `insert into currency_price(symbol) values($1)`

	_, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return err
	}

	return nil
}

func (r *CurrencyRepo) GetBySymbol(ctx context.Context, symbol string) (model.Currency, error) {
	query := "select id, symbol from currency where symbol = $1"

	row := r.db.QueryRowContext(ctx, query, &symbol)
	if row.Err() != nil {
		return model.Currency{}, row.Err()
	}

	var res model.Currency
	if err := row.Scan(
		&res.ID,
		&res.Symbol,
	); err != nil {
		return model.Currency{}, err
	}

	return res, nil
}

func (r *CurrencyRepo) List(ctx context.Context) ([]model.Currency, error) {
	query := `select id, symbol from currency`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Currency
	for rows.Next() {
		var item model.Currency
		if err := rows.Scan(
			&item.ID,
			&item.Symbol,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
