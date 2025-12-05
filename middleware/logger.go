package middleware

import (
	"time"

	httplogger "cozeva.com/vault/pkg/logger"
	"github.com/gin-gonic/gin"
)

func WinstonLogger(logger *httplogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		level := "info"
		if status >= 500 {
			level = "error"
		} else if status >= 400 {
			level = "warn"
		}

		meta := map[string]any{
			"status":   status,
			"method":   method,
			"path":     path,
			"latency":  latency.String(),
			"clientIP": clientIP,
		}

		logFuncs := map[string]func(string, map[string]any){
			"info":  logger.Info,
			"warn":  logger.Warn,
			"error": logger.Error,
		}

		if logFunc, ok := logFuncs[level]; ok {
			logFunc("HTTP Request", meta)
		} else {
			logger.Info("HTTP Request", meta)
		}
	}
}
