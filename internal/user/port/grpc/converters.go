package grpc

import (
	"time"

	"github.com/truongle2004/mercato/internal/user/model"
	pb "github.com/truongle2004/mercato/proto/gen/go/user"
)

// userInfoFromModel returns the gRPC UserInfo for a storage user. Hides internal fields
// (password, deleted_at) and removes the unreachable utils.Copy error branches.
func userInfoFromModel(m *model.User) *pb.UserInfo {
	if m == nil {
		return nil
	}
	return &pb.UserInfo{
		Id:        m.ID,
		Email:     m.Email,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
