package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	mock_repository "gexabyte/internal/repository/mock"
	mock_binance "gexabyte/pkg/clients/binance/mock"
	"strconv"
	"testing"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetStat24H(t *testing.T) {
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

	ticker24hDefaultResopnce := &binance_connector.Ticker24hrResponse{
		Symbol:    "BTCUSDT",
		OpenPrice: "1.1",
		LastPrice: "2.2",
		HighPrice: "3.3",
		LowPrice:  "4.4",
		OpenTime:  1,
		CloseTime: 2,
	}

	unexpectedErr := fmt.Errorf("unexpected")

	tc := []struct {
		name        string
		symbols     []string
		buildStubs  func(binance *mock_binance.MockClient)
		checkResult func(t *testing.T, res []model.GetCurrencyStat24HDTO, err error)
	}{
		{
			name:    "OK",
			symbols: []string{"1"},
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(ticker24hDefaultResopnce, nil)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyStat24HDTO, err error) {
				assert.NoError(t, err)

				openPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.OpenPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, openPrice, res[0].OpenPrice)
				lastPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.LastPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, lastPrice, res[0].LastPrice)
				highPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.HighPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, highPrice, res[0].HighPrice)
				lowPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.LowPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, lowPrice, res[0].LowPrice)
			},
		},
		{
			name:    "OK Multiple",
			symbols: []string{"1", "2", "3"},
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Any()).Times(3).Return(ticker24hDefaultResopnce, nil)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyStat24HDTO, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 3, len(res))

				openPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.OpenPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, openPrice, res[0].OpenPrice)
				lastPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.LastPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, lastPrice, res[0].LastPrice)
				highPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.HighPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, highPrice, res[0].HighPrice)
				lowPrice, err := strconv.ParseFloat(ticker24hDefaultResopnce.LowPrice, 64)
				assert.NoError(t, err)
				assert.Equal(t, lowPrice, res[0].LowPrice)
			},
		},
		{
			name:    "error from binance client",
			symbols: []string{"1"},
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(nil, unexpectedErr)
			},
			checkResult: func(t *testing.T, res []model.GetCurrencyStat24HDTO, err error) {
				assert.Error(t, err)
				assert.Equal(t, err, unexpectedErr)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(binanceClient)

			res, err := service.GetStat24H(context.Background(), test.symbols...)

			test.checkResult(t, res, err)
		})
	}
}

func TestFetchStat24H(t *testing.T) {
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

	ticker24hDefaultResopnce := &binance_connector.Ticker24hrResponse{
		Symbol:    "BTCUSDT",
		OpenPrice: "1.1",
		LastPrice: "2.2",
		HighPrice: "3.3",
		LowPrice:  "4.4",
		OpenTime:  1,
		CloseTime: 2,
	}

	tc := []struct {
		name        string
		symbol      string
		buildStubs  func(binance *mock_binance.MockClient)
		checkResult func(t *testing.T, res model.GetCurrencyStat24HDTO, err error)
	}{
		{
			name:   "OK",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(ticker24hDefaultResopnce, nil)
			},
			checkResult: func(t *testing.T, res model.GetCurrencyStat24HDTO, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "error parse open price",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				res := ticker24hDefaultResopnce
				res.OpenPrice = "incorrect"

				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(res, nil)
			},
			checkResult: func(t *testing.T, res model.GetCurrencyStat24HDTO, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error parse last price",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				res := ticker24hDefaultResopnce
				res.LastPrice = "incorrect"

				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(res, nil)
			},
			checkResult: func(t *testing.T, res model.GetCurrencyStat24HDTO, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error parse high price",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				res := ticker24hDefaultResopnce
				res.HighPrice = "incorrect"

				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(res, nil)
			},
			checkResult: func(t *testing.T, res model.GetCurrencyStat24HDTO, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "error parse low price",
			symbol: "1",
			buildStubs: func(binance *mock_binance.MockClient) {
				res := ticker24hDefaultResopnce
				res.LowPrice = "incorrect"

				binance.EXPECT().Ticker24hService(gomock.Any(), gomock.Eq("1")).Times(1).Return(res, nil)
			},
			checkResult: func(t *testing.T, res model.GetCurrencyStat24HDTO, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(binanceClient)

			res, err := service.fetchStat24H(context.Background(), test.symbol)

			test.checkResult(t, res, err)
		})
	}
}
