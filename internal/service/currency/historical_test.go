package currency

import (
	"context"
	"fmt"
	"gexabyte/internal/model"
	mock_repository "gexabyte/internal/repository/mock"
	mock_binance "gexabyte/pkg/clients/binance/mock"
	"testing"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSolvePagination(t *testing.T) {
	service := Currency{}

	tc := []struct {
		name string

		startTime, endTime int64
		limit, page        int
		interval           string

		checkResult func(t *testing.T, startTime int64, maxPage int)
	}{
		{
			name:      "check seconds",
			startTime: 0000,
			endTime:   5000,
			limit:     1,
			page:      1,
			interval:  "1s",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(0), startTime)
				assert.Equal(t, 6, maxPage)
			},
		},
		{
			name:      "check seconds ceil",
			startTime: 0001,
			endTime:   5000,
			limit:     1,
			page:      1,
			interval:  "1s",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1000), startTime)
				assert.Equal(t, 5, maxPage) //note cause ceil(0001) -> 1000 milliseonds
			},
		},
		{
			name:      "check seconds round",
			startTime: 0001,
			endTime:   5999, // -> 5000
			limit:     1,
			page:      1,
			interval:  "1s",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1000), startTime)
				assert.Equal(t, 5, maxPage)
			},
		},
		{
			name:      "check minute",
			startTime: 60000 * 0,
			endTime:   60000 * 5,
			limit:     1,
			page:      1,
			interval:  "1m",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(0), startTime)
				assert.Equal(t, 6, maxPage)
			},
		},
		{
			name:      "check minute ceil",
			startTime: 60000 * 1,
			endTime:   60000 * 5,
			limit:     1,
			page:      1,
			interval:  "1m",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(60*1000), startTime)
				assert.Equal(t, 5, maxPage)
			},
		},
		{
			name:      "check minute round",
			startTime: 60000 * 1,
			endTime:   60000*5 + 59*1000, // +59 sec
			limit:     1,
			page:      1,
			interval:  "1m",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(60*1000), startTime)
				assert.Equal(t, 5, maxPage)
			},
		},
		{
			name:      "check days",
			startTime: 1704067200000, // Mon 1 January 2024 00:00:00
			endTime:   1704412800000, // Fri 5 January 2024 00:00:00
			limit:     2,
			page:      1,
			interval:  "1d",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 3, maxPage)
			},
		},
		{
			name:      "check days 2 page",
			startTime: 1704067200000, // Mon 1 January 2024 00:00:00
			endTime:   1704412800000, // Fri 5 January 2024 00:00:00
			limit:     2,
			page:      2,
			interval:  "1d",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704240000000), startTime) // Mon 3 January 2024 00:00:00
				assert.Equal(t, 3, maxPage)
			},
		},
		{
			name:      "check days ceil",
			startTime: 1704067200000 + 1*1000, // Mon 1 January 2024 00:00:01 -> 2 January
			endTime:   1704412800000,          // Fri 5 January 2024 00:00:00
			limit:     2,
			page:      1,
			interval:  "1d",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704153600000), startTime) // Mon 2 January 2024 00:00:00
				assert.Equal(t, 2, maxPage)                      // note cause 4 days now
			},
		},
		{
			name:      "check week",
			startTime: 1704067200000, // Mon 1 January 	2024 00:00:00
			endTime:   1704671999999, // Sun 7 January  2024 23:59:59
			limit:     1,
			page:      1,
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 1, maxPage)
			},
		},
		{
			name:      "check week two mondays",
			startTime: 1704067200000, // Mon 1 January 	2024 00:00:00
			endTime:   1704672001000, // Mon 8 January  2024 00:00:01
			limit:     1,
			page:      1,
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 2, maxPage)
			},
		},
		{
			name:      "check week ceil",
			startTime: 1704067200001, // Mon 1 January 	2024 00:00:00 -> Mon 8 January  2024 00:00:00
			endTime:   1704672001000, // Mon 8 January  2024 00:00:01
			limit:     1,
			page:      1,
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704672000000), startTime) // Mon 8 January 2024 00:00:00
				assert.Equal(t, 1, maxPage)
			},
		},
		{
			name:      "check 2week",
			startTime: 1704067200000, // Mon 1 January 	2024 00:00:00
			endTime:   1704672000000, // Mon 8 January  2024 00:00:00
			limit:     1,
			page:      1,
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 2, maxPage)
			},
		},
		{
			name:      "check week5 limit2",
			startTime: 1704067200000, // Mon 1 January 	 2024 00:00:00
			endTime:   1706486400000, // Mon 29 January  2024 00:00:00
			limit:     2,
			page:      1,
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 3, maxPage)
			},
		},
		{
			name:      "check week5 limit2 page2",
			startTime: 1704067200000, // Mon 1 January 	 2024 00:00:00
			endTime:   1706486400000, // Mon 29 January  2024 00:00:00
			limit:     2,
			page:      2, //note
			interval:  "1w",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1705276800000), startTime) // Mon 15 January 2024 00:00:00
				assert.Equal(t, 3, maxPage)
			},
		},
		{
			name:      "check month",
			startTime: 1704067200000, // Mon 1 January 2024 00:00:00
			endTime:   1714521600000, // Wed 1 May     2024 00:00:00
			limit:     1,
			page:      1,
			interval:  "1M",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1704067200000), startTime) // Mon 1 January 2024 00:00:00
				assert.Equal(t, 5, maxPage)
			},
		},
		{
			name:      "check month ceil",
			startTime: 1704067200001, // Mon 1 January 2024 00:00:01 -> Thu 1 February 2024 00:00:00
			endTime:   1714521600000, // Wed 1 May     2024 00:00:00
			limit:     1,
			page:      1,
			interval:  "1M",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1706745600000), startTime) // Thu 1 February 2024 00:00:00
				assert.Equal(t, 4, maxPage)
			},
		},
		{
			name:      "check month div",
			startTime: 1704067200001, // Mon 1 January 2024 00:00:01 -> Thu 1 February 2024 00:00:00
			endTime:   1714521600001, // Wed 1 May     2024 00:00:01 -> Wed 1 May      2024 00:00:00
			limit:     1,
			page:      1,
			interval:  "1M",

			checkResult: func(t *testing.T, startTime int64, maxPage int) {
				assert.Equal(t, int64(1706745600000), startTime) // Thu 1 February 2024 00:00:00
				assert.Equal(t, 4, maxPage)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			st, mp := service.solvePagination(test.startTime, test.endTime, test.limit, test.page, test.interval)
			test.checkResult(t, st, mp)
		})
	}

}

func TestGetPriceHistorical(t *testing.T) {
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

	defaultReq := model.GetCurrencyPriceHistoricalDTOReq{
		Symbol:    "1",
		Interval:  "1s",
		StartTime: 1000,
		EndTime:   5000,
		Page:      1,
		Limit:     2,
	}

	defaultRes := []*binance_connector.KlinesResponse{
		{
			OpenTime:  1000,
			CloseTime: 5000,

			Open:  "1.1",
			Close: "2.2",
			High:  "3.3",
			Low:   "4.4",
		},
	}

	tc := []struct {
		name        string
		input       func() model.GetCurrencyPriceHistoricalDTOReq
		buildStubs  func(binance *mock_binance.MockClient)
		checkResult func(t *testing.T, res *model.GetCurrencyPriceHistoricalDTORes, err error)
	}{
		{
			name:  "OK",
			input: func() model.GetCurrencyPriceHistoricalDTOReq { return defaultReq },
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().KlineService(
					gomock.Any(),
					gomock.Eq("1"),
					gomock.Eq("1s"),
					gomock.Eq(int64(1000)),
					gomock.Eq(int64(5000)),
					gomock.Eq(2),
				).Times(1).Return(defaultRes, nil)
			},
			checkResult: func(t *testing.T, res *model.GetCurrencyPriceHistoricalDTORes, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, res.Page)
				assert.Equal(t, 3, res.MaxPage) // in 5sec with interval 2sec is 3 page
			},
		},
		{
			name:  "error from binance",
			input: func() model.GetCurrencyPriceHistoricalDTOReq { return defaultReq },
			buildStubs: func(binance *mock_binance.MockClient) {
				binance.EXPECT().KlineService(
					gomock.Any(),
					gomock.Eq("1"),
					gomock.Eq("1s"),
					gomock.Eq(int64(1000)),
					gomock.Eq(int64(5000)),
					gomock.Eq(2),
				).Times(1).Return(nil, unexpectedErr)
			},
			checkResult: func(t *testing.T, res *model.GetCurrencyPriceHistoricalDTORes, err error) {
				assert.Error(t, err)
				assert.Empty(t, res)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(binanceClient)

			res, err := service.GetPriceHistorical(context.Background(), test.input())

			test.checkResult(t, res, err)
		})
	}
}
