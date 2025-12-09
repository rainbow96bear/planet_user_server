package grpcclient

import (
	"github.com/rainbow96bear/planet_utils/pb"
	"google.golang.org/grpc"
)

type AuthAPI interface {
	// GetProfile(ctx context.Context, authId string) (*pb.AuthProfileResponse, error)
}

type authAPI struct {
	client pb.AuthServiceClient
}

func NewAuthAPI(conn *grpc.ClientConn) AuthAPI {
	return &authAPI{
		client: pb.NewAuthServiceClient(conn),
	}
}
