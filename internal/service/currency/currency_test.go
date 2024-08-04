package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	mock_repository "gexabyte/internal/repository/mock"
	mock_binance "gexabyte/pkg/clients/binance/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyRepo := mock_repository.NewMockCurrency(ctrl)
	currencyPriceRepo := mock_repository.NewMockCurrencyPrice(ctrl)
	binanceClient := mock_binance.NewMockClient(ctrl)

	service := Currency{
		currencyRepo:      currencyRepo,
		currencyPriceRepo: currencyPriceRepo,
		binanceClient:     binanceClient,
	}

	unexpectedErr := fmt.Errorf("unexpected")

	currencyRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	err := service.Create(context.Background(), "symbol")
	assert.NoError(t, err)

	currencyRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(unexpectedErr)
	err = service.Create(context.Background(), "symbol")
	assert.Error(t, err)
}

func TestListCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyRepo := mock_repository.NewMockCurrency(ctrl)
	currencyPriceRepo := mock_repository.NewMockCurrencyPrice(ctrl)
	binanceClient := mock_binance.NewMockClient(ctrl)

	service := Currency{
		currencyRepo:      currencyRepo,
		currencyPriceRepo: currencyPriceRepo,
		binanceClient:     binanceClient,
	}

	unexpectedErr := fmt.Errorf("unexpected")

	currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return([]model.Currency{}, nil)
	res, err := service.List(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)

	currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return(nil, unexpectedErr)
	res, err = service.List(context.Background())
	assert.Error(t, err)
	assert.Nil(t, res)
}
