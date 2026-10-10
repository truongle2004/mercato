package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/truongle2004/mercato-kit/logger"
	"go.opentelemetry.io/otel/trace"
)

const requestIDHeader = "X-Request-ID"

// RequestLogger emits one structured access log per request and correlates it
// with the active OpenTelemetry span.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Header(requestIDHeader, requestID)

		started := time.Now()
		c.Next()

		fields := []interface{}{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", time.Since(started).Milliseconds(),
		}
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		if spanContext.IsValid() {
			fields = append(fields,
				"trace_id", spanContext.TraceID().String(),
				"span_id", spanContext.SpanID().String(),
			)
		}
		logger.Infow("http request", fields...)
	}
}
