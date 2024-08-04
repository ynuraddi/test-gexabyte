package postgres

import (
	"context"
	"database/sql"
	"gexabyte/internal/model"
)

type CurrencyPriceRepo struct {
	db *sql.DB
}

func NewCurrencyPrice(db *sql.DB) *CurrencyPriceRepo {
	return &CurrencyPriceRepo{
		db: db,
	}
}

func (r *CurrencyPriceRepo) Create(ctx context.Context, rates ...model.CurrencyPrice) error {
	query := `insert into currency_price(currency_id, price, time) values($1, $2, $3)`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	for _, rate := range rates {
		_, err := tx.ExecContext(ctx, query, rate.CurrencyID, rate.Price, rate.Time)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return rbErr
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *CurrencyPriceRepo) List(ctx context.Context) ([]model.CurrencyPrice, error) {
	query := `select id, currency_id, price, time from currency_price order by time`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CurrencyPrice
	for rows.Next() {
		var item model.CurrencyPrice
		if err := rows.Scan(
			&item.ID,
			&item.CurrencyID,
			&item.Price,
			&item.Time,
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
