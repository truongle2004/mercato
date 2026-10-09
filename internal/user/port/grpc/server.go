package grpc

import (
	"github.com/truongle2004/mercato-kit/validation"
	"google.golang.org/grpc"

	"github.com/truongle2004/mercato/internal/user/repository"
	"github.com/truongle2004/mercato/internal/user/service"
	"github.com/truongle2004/mercato/pkg/dbs"
	pb "github.com/truongle2004/mercato/proto/gen/go/user"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	userRepo := repository.NewUserRepository(db)
	// gRPC stays on the local password flow; only the HTTP edge wires the
	// headless Authentik client when auth_mode=oidc.
	userSvc := service.NewUserService(validator, userRepo, nil)
	userHandler := NewUserHandler(userSvc)

	pb.RegisterUserServiceServer(svr, userHandler)
}
