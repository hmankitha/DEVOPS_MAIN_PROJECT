package middleware

import (
	"time"
	"user-management/utils"

	"github.com/gin-gonic/gin"
)

func ElasticLoggingMiddleware(elasticLogger *utils.ElasticLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Prepare structured log fields
		fields := map[string]interface{}{
			"event":        "http_request",
			"method":       method,
			"path":         path,
			"status_code":  statusCode,
			"latency_ms":   latency.Milliseconds(),
			"latency_sec":  latency.Seconds(),
			"client_ip":    clientIP,
			"user_agent":   userAgent,
			"query_params": c.Request.URL.RawQuery,
		}

		// Log based on status code
		if elasticLogger != nil {
			message := "HTTP Request"
			switch {
			case statusCode >= 500:
				elasticLogger.Error(message, fields)
			case statusCode >= 400:
				elasticLogger.Warning(message, fields)
			default:
				elasticLogger.Info(message, fields)
			}
		}
	}
}
