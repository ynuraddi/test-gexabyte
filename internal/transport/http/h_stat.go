package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStat24H godoc
//
//	@Summary		Get 24h statistics
//	@Description	Retrieves 24-hour statistics for the specified symbols.
//	@Tags			stat
//	@Produce		json
//	@Param			symbols	query		string	true	"symbols"	example(["BTCUSDT", "ETHUSDT"])
//	@Success		200		{object}	[]model.GetCurrencyStat24HDTO
//	@Failure		400		{object}	ErrMsg	"Invalid request parameters"
//	@Failure		500		{object}	ErrMsg	"Internal server error"
//	@Router			/stat/24h [get]
func (s *Server) GetStat24H(c *gin.Context) {
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

	stats, err := s.service.Currency.GetStat24H(ctx, symbols...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrMsg{err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
