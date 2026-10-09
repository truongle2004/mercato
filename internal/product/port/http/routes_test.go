package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/truongle2004/mercato-kit/validation"

	dbMocks "github.com/truongle2004/mercato/pkg/dbs/mocks"
	redisMocks "github.com/truongle2004/mercato/pkg/redis/mocks"
)

func TestRoutes(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)
	Routes(gin.New().Group("/"), mockDB, validation.New(), mockRedis)
}
