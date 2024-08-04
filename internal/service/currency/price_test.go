package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	mock_repository "gexabyte/internal/repository/mock"
	mock_binance "gexabyte/pkg/clients/binance/mock"
	"testing"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePrice(t *testing.T) {
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

	currencyPriceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	err := service.CreatePrice(context.Background(), model.CurrencyPrice{})
	assert.NoError(t, err)

	currencyPriceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(unexpectedErr)
	err = service.CreatePrice(context.Background(), model.CurrencyPrice{})
	assert.Error(t, err)
}

func TestListPrices(t *testing.T) {
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

	currencyPriceRepo.EXPECT().List(gomock.Any()).Times(1).Return([]model.CurrencyPrice{}, nil)
	res, err := service.ListPrices(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)

	currencyPriceRepo.EXPECT().List(gomock.Any()).Times(1).Return(nil, unexpectedErr)
	res, err = service.ListPrices(context.Background())
	assert.Error(t, err)
	assert.Nil(t, res)
}

type currencyPriceMatcher struct {
	currencyIDPrice map[int]float64
}

func (c currencyPriceMatcher) Matches(x interface{}) bool {
	prices, ok := x.([]model.CurrencyPrice)
	if !ok {
		return ok
	}

	for _, p := range prices {
		v, ok := c.currencyIDPrice[p.CurrencyID]
		if !ok || v != p.Price {
			return false
		}
	}

	return true
}

func (c currencyPriceMatcher) String() string {
	return "metches []model.CurrencyPrice"
}

func TestGetCurrentPrices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyRepo := mock_repository.NewMockCurrency(ctrl)
	currencyPriceRepo := mock_repository.NewMockCurrencyPrice(ctrl)
	binanceClient := mock_binance.NewMockClient(ctrl)

	service := Currency{
		currencyRepo:      currencyRepo,
		currencyPriceRepo: currencyPriceRepo,
		binanceClient:     binanceClient,

		priceCheckTicker:   time.NewTicker(10 * time.Minute),
		priceCheckInterval: 10 * time.Minute,
	}

	unexpectedErr := fmt.Errorf("unexpected")

	tc := []struct {
		name        string
		symbols     []string
		buildStubs  func()
		checkResult func(t *testing.T, res []model.GetCurrencyPriceDTO, err error)
	}{
		{
			name:    "OK",
			symbols: []string{"1"},
			buildStubs: func() {
				currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return([]model.Currency{{ID: 1, Symbol: "1"}}, nil)

				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Eq("1")).Times(1).Return(&binance_connector.TickerPriceResponse{Symbol: "1", Price: "1.1"}, nil)

				currencyPriceRepo.EXPECT().Create(gomock.Any(), currencyPriceMatcher{
					currencyIDPrice: map[int]float64{
						1: 1.1,
					},
				}).Times(1).Return(nil)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyPriceDTO, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, res[0].Symbol, "1")
				assert.Equal(t, res[0].Price, 1.1)
			},
		},
		{
			name:    "OK no tracked symbol",
			symbols: []string{"2"},
			buildStubs: func() {
				// note than we have tracked symbol in bd, that will refresh, but will not return
				currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return([]model.Currency{{ID: 1, Symbol: "1"}}, nil)

				// mb flucky test cause here is possible race
				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Any()).Times(1).Return(&binance_connector.TickerPriceResponse{Symbol: "2", Price: "2.2"}, nil)
				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Any()).Times(1).Return(&binance_connector.TickerPriceResponse{Symbol: "1", Price: "1.1"}, nil)

				currencyPriceRepo.EXPECT().Create(gomock.Any(), currencyPriceMatcher{
					currencyIDPrice: map[int]float64{
						1: 1.1,
					},
				}).Times(1).Return(nil)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyPriceDTO, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)

				expectedSymbolPrice := map[string]float64{
					"2": 2.2,
				}

				for _, r := range res {
					v, ok := expectedSymbolPrice[r.Symbol]
					assert.True(t, ok)
					assert.Equal(t, r.Price, v)
				}

			},
		},
		{
			name:    "OK currency table is empty",
			symbols: []string{"2"},
			buildStubs: func() {
				currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return(nil, nil)
				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Any()).Times(1).Return(&binance_connector.TickerPriceResponse{Symbol: "2", Price: "2.2"}, nil)
				currencyPriceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyPriceDTO, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)

				expectedSymbolPrice := map[string]float64{
					"2": 2.2,
				}

				for _, r := range res {
					v, ok := expectedSymbolPrice[r.Symbol]
					assert.True(t, ok)
					assert.Equal(t, r.Price, v)
				}

			},
		},
		{
			name:    "error from currency db",
			symbols: []string{"1"},
			buildStubs: func() {
				currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return(nil, unexpectedErr)
				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Any()).Times(0)
				currencyPriceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyPriceDTO, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name:    "error from ticker service",
			symbols: []string{"1"},
			buildStubs: func() {
				currencyRepo.EXPECT().List(gomock.Any()).Times(1).Return(nil, nil)
				binanceClient.EXPECT().TickerPriceService(gomock.Any(), gomock.Any()).Times(1).Return(nil, unexpectedErr)
				currencyPriceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyPriceDTO, err error) {
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs()

			res, err := service.GetCurrentPrices(context.Background(), test.symbols...)

			test.checkResult(t, res, err)
		})
	}
}

func TestFetchCurrentPrice(t *testing.T) {
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

	defaultCurrentPrice := &binance_connector.TickerPriceResponse{
		Symbol: "1",
		Price:  "1",
	}

	tc := []struct {
		name        string
		symbol      string
		buildStubs  func(binance *mock_binance.MockClient)
		checkResult func(t *testing.T, res float64, err error)
	}{
		{
			name:   "OK",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().TickerPriceService(gomock.Any(), gomock.Eq("1")).Times(1).Return(defaultCurrentPrice, nil)
			},
			checkResult: func(t *testing.T, res float64, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error parse float",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				cp := *defaultCurrentPrice
				cp.Price = "incorrect"
				binance.EXPECT().TickerPriceService(gomock.Any(), gomock.Eq("1")).Times(1).Return(&cp, nil)
			},
			checkResult: func(t *testing.T, res float64, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error from binance",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().TickerPriceService(gomock.Any(), gomock.Eq("1")).Times(1).Return(nil, unexpectedErr)
			},
			checkResult: func(t *testing.T, res float64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(binanceClient)

			res, err := service.fetchCurrentPrice(context.Background(), test.symbol)

			test.checkResult(t, res, err)
		})
	}
}
