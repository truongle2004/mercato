package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/truongle2004/mercato-kit/logger"
)

func TestRequestLogger_AddsRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Initialize("test")

	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/products", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusNoContent, response.Code)
	require.NotEmpty(t, response.Header().Get(requestIDHeader))
}

func TestRequestLogger_PreservesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Initialize("test")

	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/products", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	request.Header.Set(requestIDHeader, "request-123")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, "request-123", response.Header().Get(requestIDHeader))
}
