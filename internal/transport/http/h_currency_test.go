package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gexabyte/internal/model"
	"gexabyte/internal/service"
	mock_service "gexabyte/internal/service/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCurrency(t *testing.T) {
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
		symbol        string
		buildStubs    func(service *mock_service.MockCurrency)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			symbol: "BTCUSDT",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().Create(gomock.Any(), gomock.Eq("BTCUSDT")).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:   "empty field",
			symbol: "",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "lower case field",
			symbol: "bTCUSDT",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "internal server error",
			symbol: "BTCUSDT",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("unexpected"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService)

			reqBody := CreateCurrencyReq{Symbol: test.symbol}
			body, err := json.Marshal(reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/currency", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			router := server.setupRouter()
			router.ServeHTTP(rec, req)

			test.checkResponse(t, rec)
		})
	}

}

func TestListCurrency(t *testing.T) {
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
				service.EXPECT().List(gomock.Any()).Times(1).Return([]model.Currency{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "internal server error",
			buildStubs: func(service *mock_service.MockCurrency) {
				service.EXPECT().List(gomock.Any()).Times(1).Return(nil, fmt.Errorf("unexpected"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies", nil)
			rec := httptest.NewRecorder()

			router := server.setupRouter()
			router.ServeHTTP(rec, req)

			test.checkResponse(t, rec)
		})
	}

}
