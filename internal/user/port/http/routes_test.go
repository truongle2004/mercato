package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/truongle2004/mercato-kit/validation"

	"github.com/truongle2004/mercato/pkg/dbs/mocks"
)

func TestRoutes(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	Routes(gin.New().Group("/"), mockDB, validation.New())
}
