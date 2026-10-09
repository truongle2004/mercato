package grpc

import (
	"github.com/truongle2004/mercato-kit/validation"
	"google.golang.org/grpc"

	notificationRepo "github.com/truongle2004/mercato/internal/notification/repository"
	notificationSvc "github.com/truongle2004/mercato/internal/notification/service"
	"github.com/truongle2004/mercato/internal/order/repository"
	"github.com/truongle2004/mercato/internal/order/service"
	userRepo "github.com/truongle2004/mercato/internal/user/repository"
	"github.com/truongle2004/mercato/pkg/config"
	"github.com/truongle2004/mercato/pkg/dbs"
	"github.com/truongle2004/mercato/pkg/notification"
	pb "github.com/truongle2004/mercato/proto/gen/go/order"
)

func RegisterHandlers(svr *grpc.Server, db dbs.Database, validator validation.Validation) {
	oRepo := repository.NewOrderRepository(db)
	pRepo := repository.NewProductRepository(db)
	uRepo := repository.NewUserRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	cfg := config.GetConfig()
	couponSvc := service.NewCouponService(validator, couponRepo)
	prefChecker := notificationSvc.NewDBPreferenceChecker(
		notificationSvc.NewUserRepoLookup(userRepo.NewUserRepository(db)),
		notificationRepo.NewPreferenceRepository(db),
	)
	notifier := notification.BuildDefault(notification.Settings{
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUser:     cfg.SMTPUser,
		SMTPPassword: cfg.SMTPPassword,
		EmailFrom:    cfg.EmailFrom,
		Prefs:        prefChecker,
		DLQ:          notificationRepo.NewDeadLetterSink(db),
	})
	orderSvc := service.NewOrderService(validator, db, oRepo, pRepo, uRepo, reservationRepo, couponSvc, notifier)
	orderHandler := NewOrderHandler(orderSvc)

	pb.RegisterOrderServiceServer(svr, orderHandler)
}
