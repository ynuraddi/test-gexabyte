package http

import (
	"fmt"
	"gexabyte/internal/model"
	"gexabyte/internal/service"
	mock_service "gexabyte/internal/service/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListPrices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyService := mock_service.NewMockCurrency(ctrl)
	service := service.Manager{Currency: currencyService}

	server := Server{
		service: &service,
		logger:  slog.Default(),
	}

	tc := []struct {
		name          string
		buildStubs    func(service *mock_service.MockCurrency)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().ListPrices(gomock.Any()).Times(1).Return([]model.CurrencyPrice{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "internal server error",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().ListPrices(gomock.Any()).Times(1).Return(nil, fmt.Errorf("unexpected"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/prices", nil)
			rec := httptest.NewRecorder()

			router := server.setupRouter()
			router.ServeHTTP(rec, req)

			test.checkResponse(t, rec)
		})
	}
}

func TestListPricesHistorical(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyService := mock_service.NewMockCurrency(ctrl)
	service := service.Manager{Currency: currencyService}

	server := Server{
		service: &service,
		logger:  slog.Default(),
	}

	sT := time.Now()
	eT := sT.Add(5 * time.Second)

	defaultParam := model.GetCurrencyPriceHistoricalDTOReq{
		Symbol:    "BTCUSDT",
		Interval:  "1s",
		StartTime: sT.UnixMilli(),
		EndTime:   eT.UnixMilli(),
		Page:      1,
		Limit:     2,
	}

	tc := []struct {
		name          string
		params        func() model.GetCurrencyPriceHistoricalDTOReq
		buildStubs    func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "empty symbol",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				defaultParam.Symbol = ""
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "empty interval",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				defaultParam.Interval = ""
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "interval incorrect format",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				defaultParam.Interval = "2s"
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "endTime earlier than startTime",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				now := time.Now()
				defaultParam.EndTime = now.UnixMilli()
				defaultParam.StartTime = now.Add(5 * time.Second).UnixMilli()
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "negative limit",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				defaultParam.Limit = -1
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "negative page",
			params: func() model.GetCurrencyPriceHistoricalDTOReq {
				defaultParam.Page = -1
				return defaultParam
			},
			buildStubs: func(service *mock_service.MockCurrency, eq model.GetCurrencyPriceHistoricalDTOReq) {
				service.EXPECT().GetPriceHistorical(gomock.Any(), gomock.Eq(eq)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService, test.params())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/prices/historical", nil)

			q := req.URL.Query()
			if test.params().Symbol != "" {
				q.Add("symbol", test.params().Symbol)
			}
			if test.params().Interval != "" {
				q.Add("interval", test.params().Interval)
			}
			if test.params().StartTime != 0 {
				q.Add("startTime", strconv.Itoa(int(test.params().StartTime)))
			}
			if test.params().EndTime != 0 {
				q.Add("endTime", strconv.Itoa(int(test.params().EndTime)))
			}
			if test.params().Page != 0 {
				q.Add("page", strconv.Itoa(int(test.params().Page)))
			}
			if test.params().Limit != 0 {
				q.Add("limit", strconv.Itoa(int(test.params().Limit)))
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			router := server.setupRouter()
			router.ServeHTTP(rec, req)

			test.checkResponse(t, rec)
		})
	}
}

func TestListPricesCurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currencyService := mock_service.NewMockCurrency(ctrl)
	service := service.Manager{Currency: currencyService}

	server := Server{
		service: &service,
		logger:  slog.Default(),
	}

	tc := []struct {
		name          string
		query         string
		value         string
		buildStubs    func(service *mock_service.MockCurrency)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "symbols",
			value: `["BTCUSDT"]`,
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().GetCurrentPrices(gomock.Any(), gomock.Eq([]string{"BTCUSDT"})).Times(1).Return([]model.GetCurrencyPriceDTO{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "OK",
			query: "symbols",
			value: `["BTCUSDT", "ETHUSDT"]`,
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().GetCurrentPrices(gomock.Any(), gomock.Eq([]string{"BTCUSDT", "ETHUSDT"})).Times(1).Return([]model.GetCurrencyPriceDTO{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "bad request symbols param is required",
			query: "",
			value: `["BTCUSDT"]`,
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().GetCurrentPrices(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "bad request invalid symbol format",
			query: "symbols",
			value: `BTCUSDT"]`,
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().GetCurrentPrices(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "internal server error",
			query: "symbols",
			value: `["BTCUSDT"]`,
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().GetCurrentPrices(gomock.Any(), gomock.Any()).Times(1).Return(nil, fmt.Errorf("unexpected"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/prices/current", nil)

			if test.query != "" {
				q := req.URL.Query()
				q.Add(test.query, test.value)
				req.URL.RawQuery = q.Encode()
			}

			rec := httptest.NewRecorder()

			router := server.setupRouter()
			router.ServeHTTP(rec, req)

			test.checkResponse(t, rec)
		})
	}
}
