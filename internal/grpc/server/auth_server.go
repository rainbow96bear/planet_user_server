package grpcserver

import (
	"net"

	"github.com/rainbow96bear/planet_user_server/config"
	grpcclient "github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	pb "github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"google.golang.org/grpc"
)

type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
	Clients *grpcclient.GrpcClients
}

func NewUserGrpcServer(clients *grpcclient.GrpcClients) *UserGrpcServer {
	return &UserGrpcServer{
		Clients: clients,
	}
}

func RunGrpcServer() error {
	// 🔥 1) 모든 gRPC 클라이언트 생성
	clients, err := grpcclient.NewGrpcClients()
	if err != nil {
		return err
	}

	// 🔥 2) gRPC 서버 Listen 시작
	listener, err := net.Listen("tcp", ":"+config.USER_GRPC_PORT)
	if err != nil {
		return err
	}

	// 🔥 3) gRPC 서버 생성
	grpcServer := grpc.NewServer()

	// 🔥 4) UserGrpcServer 등록
	authServer := NewUserGrpcServer(clients)
	pb.RegisterUserServiceServer(grpcServer, authServer)

	logger.Debugf("User gRPC Server running on :%s\n", config.USER_GRPC_PORT)

	// 🔥 5) 서버 시작
	return grpcServer.Serve(listener)
}
