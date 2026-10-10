package grpc

import (
	"fmt"
	"net"

	"github.com/truongle2004/mercato-kit/logger"
	"github.com/truongle2004/mercato-kit/validation"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	orderGRPC "github.com/truongle2004/mercato/internal/order/port/grpc"
	productGRPC "github.com/truongle2004/mercato/internal/product/port/grpc"
	userGRPC "github.com/truongle2004/mercato/internal/user/port/grpc"
	"github.com/truongle2004/mercato/pkg/config"
	"github.com/truongle2004/mercato/pkg/dbs"
	"github.com/truongle2004/mercato/pkg/middleware"
	"github.com/truongle2004/mercato/pkg/oidc"
	"github.com/truongle2004/mercato/pkg/redis"
)

type Server struct {
	engine    *grpc.Server
	cfg       *config.Schema
	validator validation.Validation
	db        dbs.Database
	cache     redis.Redis
}

func NewServer(validator validation.Validation, db dbs.Database, cache redis.Redis) *Server {
	cfg := config.GetConfig()
	var interceptor grpc.UnaryServerInterceptor

	switch cfg.AuthMode {
	case config.AuthModeOIDC:
		oidcValidator, err := oidc.NewValidator(cfg.OIDCIssuer, cfg.OIDCClientID, cfg.OIDCClientSecret, cfg.OIDCRedirectURL)
		if err != nil {
			logger.Fatalf("Failed to create OIDC validator: %v", err)
		}
		interceptor = middleware.NewOIDCInterceptor(oidcValidator, config.AuthIgnoreMethods).Unary()
	default:
		interceptor = middleware.NewAuthInterceptor(config.AuthIgnoreMethods).Unary()
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor,
		),
	)

	return &Server{
		engine:    grpcServer,
		cfg:       cfg,
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

func (s *Server) Run() error {
	userGRPC.RegisterHandlers(s.engine, s.db, s.validator)
	productGRPC.RegisterHandlers(s.engine, s.db, s.validator)
	orderGRPC.RegisterHandlers(s.engine, s.db, s.validator)

	reflection.Register(s.engine)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
	logger.Info("GRPC server is listening on PORT: ", s.cfg.GrpcPort)
	if err != nil {
		logger.Error("Failed to listen: ", err)
		return err
	}

	// Start grpc server
	err = s.engine.Serve(lis)
	if err != nil {
		logger.Fatal("Failed to serve grpc: ", err)
		return err
	}

	return nil
}

func (s *Server) Shutdown() {
	logger.Info("Shutting down gRPC server...")
	s.engine.GracefulStop()
}
