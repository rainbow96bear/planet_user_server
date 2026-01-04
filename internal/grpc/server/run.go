package server

import (
	"net"

	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_user_server/internal/bootstrap"
	pb "github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"google.golang.org/grpc"
)

func RunGrpcServer(deps *bootstrap.Dependencies) error {
	listener, err := net.Listen("tcp", ":"+config.USER_GRPC_PORT)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	userServer := NewUserGrpcServer(deps.Services.Profile)
	pb.RegisterUserServiceServer(grpcServer, userServer)

	logger.Infof(
		"User gRPC Server running on :%s",
		config.USER_GRPC_PORT,
	)

	return grpcServer.Serve(listener)
}
