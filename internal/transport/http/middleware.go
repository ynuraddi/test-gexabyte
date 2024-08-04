package http

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()

		if status >= 400 {
			s.logger.Error(fmt.Sprintf("method:%s path:%s status:%d latency:%v", c.Request.Method, c.Request.URL.Path, status, latency))
		} else {
			s.logger.Info(fmt.Sprintf("method:%s path:%s status:%d latency:%v", c.Request.Method, c.Request.URL.Path, status, latency))
		}
	}
}
