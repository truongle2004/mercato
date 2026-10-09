package grpc

import (
	"testing"

	"github.com/truongle2004/mercato-kit/validation"
	goGRPC "google.golang.org/grpc"

	"github.com/truongle2004/mercato/pkg/dbs/mocks"
)

func TestRegisterHandlers(t *testing.T) {
	mockDB := mocks.NewDatabase(t)
	RegisterHandlers(goGRPC.NewServer(), mockDB, validation.New())
}
