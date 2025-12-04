package grpcclient

import (
	"github.com/rainbow96bear/planet_utils/pb"
	"google.golang.org/grpc"
)

type UserAPI interface {
	// GetProfile(ctx context.Context, userId string) (*pb.UserProfileResponse, error)
}

type userAPI struct {
	client pb.UserServiceClient
}

func NewUserAPI(conn *grpc.ClientConn) UserAPI {
	return &userAPI{
		client: pb.NewUserServiceClient(conn),
	}
}

// func (u *userAPI) GetProfile(ctx context.Context, userId string) (*pb.UserProfileResponse, error) {
// 	return u.client.GetUserProfile(ctx, &pb.UserProfileRequest{
// 		UserId: userId,
// 	})
// }
