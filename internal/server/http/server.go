package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/truongle2004/mercato-kit/logger"
	"github.com/truongle2004/mercato-kit/validation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "github.com/truongle2004/mercato/docs"
	notificationHttp "github.com/truongle2004/mercato/internal/notification/port/http"
	orderHttp "github.com/truongle2004/mercato/internal/order/port/http"
	paymentHttp "github.com/truongle2004/mercato/internal/payment/port/http"
	productHttp "github.com/truongle2004/mercato/internal/product/port/http"
	userHttp "github.com/truongle2004/mercato/internal/user/port/http"
	"github.com/truongle2004/mercato/pkg/config"
	"github.com/truongle2004/mercato/pkg/dbs"
	"github.com/truongle2004/mercato/pkg/middleware"
	"github.com/truongle2004/mercato/pkg/redis"
	"github.com/truongle2004/mercato/pkg/response"
)

type Server struct {
	engine    *gin.Engine
	httpSvr   *http.Server
	cfg       *config.Schema
	validator validation.Validation
	db        dbs.Database
	cache     redis.Redis
}

func NewServer(validator validation.Validation, db dbs.Database, cache redis.Redis) *Server {
	return &Server{
		engine:    gin.Default(),
		cfg:       config.GetConfig(),
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

func (s *Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine.Use(otelgin.Middleware(s.cfg.TelemetryServiceName))
	s.engine.Use(middleware.RequestLogger())
	s.engine.Use(middleware.CORS())
	s.engine.Use(middleware.RateLimit(s.cache))

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, nil)
	})

	s.httpSvr = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.HttpPort),
		Handler:           s.engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start http server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	if err := s.httpSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.httpSvr.Shutdown(ctx)
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	userHttp.Routes(v1, s.db, s.validator)
	productHttp.Routes(v1, s.db, s.validator, s.cache)
	orderHttp.Routes(v1, s.db, s.validator)
	paymentHttp.Routes(v1, s.db, s.validator)
	notificationHttp.Routes(v1, s.db)
	return nil
}
