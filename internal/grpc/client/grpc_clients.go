package grpcclient

import (
	"github.com/rainbow96bear/planet_auth_server/config"
	"google.golang.org/grpc"
)

type GrpcClients struct {
	User UserAPI
	// 앞으로 증가할 클라이언트들
	// FeedClient pb.FeedServiceClient
	// ChatClient pb.ChatServiceClient
}

func NewGrpcClients() (*GrpcClients, error) {
	userConn, err := grpc.Dial(config.USER_GRPC_SERVER_ADDR, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		User: NewUserAPI(userConn),
	}, nil
}
