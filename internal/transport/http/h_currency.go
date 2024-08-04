package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCurrencyReq struct {
	Symbol string `json:"symbol" binding:"required,uppercase"`
}

// CreateCurrency godoc
//
//	@Summary		Create
//	@Description	Creates a new tracked pair.
//	@Tags			currency
//	@Accept			json
//	@Produce		json
//	@Param			currency	body	CreateCurrencyReq	true	"Currency to create"
//	@Success		201
//	@Failure		400	{object}	ErrMsg	"Invalid request parameters"
//	@Failure		500	{object}	ErrMsg	"Internal server error"
//	@Router			/currency [post]
func (s *Server) CreateCurrency(c *gin.Context) {
	var req CreateCurrencyReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrMsg{err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Copy(), 5*time.Second)
	defer cancel()

	if err := s.service.Currency.Create(ctx, req.Symbol); err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// ListPrices godoc
//
//	@Summary		List currencies
//	@Description	Retrieves a list of tracked currencies.
//	@Tags			currency
//	@Produce		json
//	@Success		200	{array}		[]model.Currency	"Retrieves a list of tracked currencies"
//	@Failure		500	{object}	ErrMsg				"Internal server error"
//	@Router			/currencies [get]
func (s *Server) ListCurrencies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Copy(), 5*time.Second)
	defer cancel()

	res, err := s.service.Currency.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
