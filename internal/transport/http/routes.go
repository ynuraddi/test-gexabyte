package http

import (
	"gexabyte/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter() *gin.Engine {
	r := gin.New()
	s.Handler = r

	r.Use(gin.Recovery(), s.LoggerMiddleware())
	api := r.Group("/api/v1")

	api.GET("/ping", s.ping)

	api.POST("/currency", s.CreateCurrency)
	api.GET("/currencies", s.ListCurrencies)

	api.GET("/prices", s.ListPrices)
	api.GET("/prices/current", s.ListPricesCurrent)
	api.GET("/prices/historical", s.ListPricesHistorical)

	api.GET("/stat/24h", s.GetStat24H)

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

// PingExample godoc
//
//	@Summary		Ping endpoint
//	@Description	Returns a 200 OK status to indicate the service is up and running
//	@Tags			ping
//	@Success		200
//	@Router			/ping [get]
func (s *Server) ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
