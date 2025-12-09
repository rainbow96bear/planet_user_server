package grpcserver

import (
	"context"
	"net"

	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/bootstrap"
	grpcclient "github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	pb "github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
	Clients        *grpcclient.GrpcClients
	ProfileService service.ProfileServiceInterface
}

func NewUserGrpcServer(
	clients *grpcclient.GrpcClients,
	profileSvc service.ProfileServiceInterface,
) *UserGrpcServer {
	return &UserGrpcServer{
		Clients:        clients,
		ProfileService: profileSvc,
	}
}

func RunGrpcServer(db *gorm.DB, deps *bootstrap.Dependencies) error {
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
	userServer := NewUserGrpcServer(clients, deps.Services.Profile)
	pb.RegisterUserServiceServer(grpcServer, userServer)

	logger.Debugf("User gRPC Server running on :%s\n", config.USER_GRPC_PORT)

	// 🔥 5) 서버 시작
	return grpcServer.Serve(listener)
}

func (s *UserGrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	logger.Debugf("Received CreateUser request: userId=%s, nickname=%s", req.UserId, req.Nickname)

	// 2. Profile 구조체 생성
	profile, err := dto.FromGrpcCreateUserRequest(req)
	if err != nil {
		return &pb.CreateUserResponse{
			Success: false,
			Message: "invalid userId",
		}, nil
	}
	// 3. DB에 저장
	_, err = s.ProfileService.CreateProfile(ctx, profile)
	if err != nil {
		logger.Errorf("Failed to create profile: %v", err)
		return &pb.CreateUserResponse{
			Success: false,
			Message: "failed to create profile",
		}, nil
	}

	// 4. 성공 응답
	return &pb.CreateUserResponse{
		Success: true,
		Message: "profile created successfully",
	}, nil
}
