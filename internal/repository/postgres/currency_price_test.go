package postgres

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyPrice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repo := NewCurrencyPrice(db)

	now := time.Now().UnixMilli()
	in := []model.CurrencyPrice{}
	in = append(in,
		model.CurrencyPrice{
			CurrencyID: 1,
			Price:      1,
			Time:       now,
		},
		model.CurrencyPrice{
			CurrencyID: 1,
			Price:      2,
			Time:       now,
		},
	)

	mock.ExpectBegin()
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(1), now).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(2), now).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	assert.NoError(t, repo.Create(context.Background(), in...))

	expectedErr := fmt.Errorf("some error")

	mock.ExpectBegin()
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(1), now).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(2), now).WillReturnError(expectedErr)
	mock.ExpectRollback()
	assert.Error(t, expectedErr, repo.Create(context.Background(), in...))

	mock.ExpectBegin()
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(1), now).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into currency_price").WithArgs(1, float64(2), now).WillReturnError(fmt.Errorf("other error"))
	mock.ExpectRollback().WillReturnError(expectedErr)
	assert.Error(t, expectedErr, repo.Create(context.Background(), in...))
}
