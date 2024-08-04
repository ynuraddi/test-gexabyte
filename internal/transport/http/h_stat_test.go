package http

import (
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

func TestGetStat24H(t *testing.T) {
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
				service.EXPECT().GetStat24H(gomock.Any(), gomock.Eq([]string{"BTCUSDT"})).Times(1).Return([]model.GetCurrencyStat24HDTO{}, nil)
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
				service.EXPECT().GetStat24H(gomock.Any(), gomock.Eq([]string{"BTCUSDT", "ETHUSDT"})).Times(1).Return([]model.GetCurrencyStat24HDTO{}, nil)
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
				service.EXPECT().GetStat24H(gomock.Any(), gomock.Any()).Times(0)
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
				service.EXPECT().GetStat24H(gomock.Any(), gomock.Any()).Times(0)
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
				service.EXPECT().GetStat24H(gomock.Any(), gomock.Any()).Times(1).Return(nil, fmt.Errorf("unexpected"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range tc {
		t.Run(test.name, func(t *testing.T) {
			test.buildStubs(currencyService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/stat/24h", nil)

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
