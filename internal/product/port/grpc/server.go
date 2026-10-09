package grpc

import (
	"github.com/truongle2004/mercato-kit/validation"
	"google.golang.org/grpc"

	"github.com/truongle2004/mercato/internal/product/repository"
	"github.com/truongle2004/mercato/internal/product/service"
	"github.com/truongle2004/mercato/pkg/dbs"
	pb "github.com/truongle2004/mercato/proto/gen/go/product"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(validator, productRepo)
	productHandler := NewProductHandler(productSvc)

	pb.RegisterProductServiceServer(svr, productHandler)
}
