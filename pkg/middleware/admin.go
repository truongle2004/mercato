package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/truongle2004/mercato/pkg/apperror"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			apperror.ErrForbidden.HTTPError(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
