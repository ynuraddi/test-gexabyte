package postgres

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	repo := NewCurrency(db)

	symbol := "BTCUSDT"

	mock.ExpectExec("insert into currency").WithArgs(symbol).WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, repo.Create(context.Background(), symbol))

	mock.ExpectExec("insert into currency").WithArgs(symbol).WillReturnError(fmt.Errorf("duplicate value"))
	assert.Error(t, repo.Create(context.Background(), symbol))

	mock.ExpectQuery("select id, symbol from currency").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "symbol"}).AddRow(1, symbol))
	res, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, model.Currency{ID: 1, Symbol: symbol}, res[0])
}
