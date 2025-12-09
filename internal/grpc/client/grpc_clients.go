package grpcclient

import (
	"github.com/rainbow96bear/planet_user_server/config"
	"google.golang.org/grpc"
)

type GrpcClients struct {
	Auth AuthAPI
}

func NewGrpcClients() (*GrpcClients, error) {
	authConn, err := grpc.Dial(config.AUTH_GRPC_SERVER_ADDR, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		Auth: NewAuthAPI(authConn),
	}, nil
}
