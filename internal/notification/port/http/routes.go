package http

import (
	"github.com/gin-gonic/gin"

	"github.com/truongle2004/mercato/internal/notification/repository"
	"github.com/truongle2004/mercato/internal/notification/service"
	"github.com/truongle2004/mercato/pkg/dbs"
	"github.com/truongle2004/mercato/pkg/middleware"
)

func Routes(r *gin.RouterGroup, db dbs.Database) {
	repo := repository.NewPreferenceRepository(db)
	svc := service.NewPreferenceService(repo)
	h := NewHandler(svc)

	g := r.Group("/me/notification-preferences", middleware.JWTAuth())
	g.GET("", h.ListPreferences)
	g.PUT("", h.SetPreference)
}
