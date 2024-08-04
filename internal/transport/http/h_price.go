package http

import (
	"context"
	"encoding/json"
	"gexabyte/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ListPrices godoc
//
//	@Summary		List currency prices
//	@Description	Retrieves a list of current currency prices.
//	@Tags			prices
//	@Produce		json
//	@Success		200	{array}		[]model.CurrencyPrice	"A list of current currency prices"
//	@Failure		500	{object}	ErrMsg					"Internal server error"
//	@Router			/prices [get]
func (s *Server) ListPrices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Copy(), 5*time.Second)
	defer cancel()

	res, err := s.service.Currency.ListPrices(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListPricesCurrent godoc
//
//	@Summary		Get purrent prices of symbols
//	@Description	Retrieves current prices fof symbols and save it in db.
//	@Tags			prices
//	@Produce		json
//	@Param			symbols	query		string	true	"symbols"	example(["BTCUSDT", "ETHUSDT"])
//	@Success		200		{object}	[]model.GetCurrencyPriceDTO
//	@Failure		400		{object}	ErrMsg	"Invalid request parameters"
//	@Failure		500		{object}	ErrMsg	"Internal server error"
//	@Router			/prices/current [get]
func (s *Server) ListPricesCurrent(c *gin.Context) {
	symbolsParam := c.Query("symbols")
	if len(symbolsParam) == 0 {
		c.JSON(http.StatusBadRequest, ErrMsg{"symbols param is required"})
		return
	}

	var symbols []string
	if err := json.Unmarshal([]byte(symbolsParam), &symbols); err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{"invalid symbols format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Copy(), 5*time.Second)
	defer cancel()

	stats, err := s.service.Currency.GetCurrentPrices(ctx, symbols...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ListPricesHistorical godoc
//
//	@Summary		List historical currency prices
//	@Description	Retrieves historical prices for a currency based on the specified parameters. Requires `symbol`, `interval`, `startTime`, `endTime`, `page`, and `limit` query parameters.
//	@Tags			prices
//	@Produce		json
//	@Param			symbol		query		string										true	"Currency symbol"
//	@Param			interval	query		string										true	"Interval"	Enums(1s, 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 1w, 1M)
//	@Param			startTime	query		int64										true	"Start time in Unix timestamp milliseconds"
//	@Param			endTime		query		int64										true	"End time in Unix timestamp milliseconds"
//	@Param			page		query		int											true	"Page number"		minimum(1)
//	@Param			limit		query		int											true	"Max limit is 1000"	minimum(1)	maximum(1000)
//	@Success		200			{object}	[]model.GetCurrencyPriceHistoricalDTORes	"Successful response with historical price data"
//	@Failure		400			{object}	ErrMsg										"Invalid request parameters"
//	@Failure		500			{object}	ErrMsg										"Internal server error"
//	@Router			/prices/historical [get]
func (s *Server) ListPricesHistorical(c *gin.Context) {
	var req model.GetCurrencyPriceHistoricalDTOReq

	req.Symbol = c.Query("symbol")
	req.Interval = c.Query("interval")

	sT, err := strconv.Atoi(c.Query("startTime"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{err.Error()})
		return
	}
	req.StartTime = int64(sT)

	eT, err := strconv.Atoi(c.Query("endTime"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{err.Error()})
		return
	}
	req.EndTime = int64(eT)

	if time.UnixMilli(req.EndTime).Before(time.UnixMilli(req.StartTime)) {
		c.JSON(http.StatusBadRequest, ErrMsg{"incorrect time - startTime is later than endTime"})
		return
	}

	req.Page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{err.Error()})
		return
	}

	req.Limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{err.Error()})
		return
	}

	if req.Symbol == "" || req.Interval == "" ||
		req.StartTime <= 0 || req.EndTime <= 0 ||
		req.Limit <= 0 || req.Page <= 0 {
		c.JSON(http.StatusBadRequest, ErrMsg{"all params are required"})
		return
	}

	if !model.KlineInterval.IsCorrect(req.Interval) {
		c.JSON(http.StatusBadRequest, ErrMsg{"incorrect interval format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Copy(), 5*time.Second)
	defer cancel()

	result, err := s.service.Currency.GetPriceHistorical(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
