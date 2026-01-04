package server

import (
	"context"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	pb "github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

// UserGrpcServer는 UserService gRPC 서버 구현체
type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
	profileService service.ProfileServiceInterface
}

// NewUserGrpcServer는 UserGrpcServer 생성자
func NewUserGrpcServer(
	profileSvc service.ProfileServiceInterface,
) *UserGrpcServer {
	return &UserGrpcServer{
		profileService: profileSvc,
	}
}

// ✅ proto에 정의된 rpc CreateProfile 구현
func (s *UserGrpcServer) CreateProfile(
	ctx context.Context,
	req *pb.CreateProfileRequest,
) (*pb.CreateProfileResponse, error) {

	logger.Infof(
		"grpc CreateProfile userId=%s nickname=%s",
		req.UserId,
		req.Nickname,
	)

	// gRPC 요청 → 도메인 DTO 변환
	profile, err := dto.FromGrpcCreateUserRequest(req)
	if err != nil {
		logger.Warnf(
			"invalid CreateProfile request userId=%s err=%v",
			req.UserId,
			err,
		)
		return &pb.CreateProfileResponse{
			Success: false,
			Message: "invalid request",
		}, nil
	}

	// 프로필 생성
	if _, err := s.profileService.CreateProfile(ctx, profile); err != nil {
		logger.Errorf(
			"CreateProfile failed userId=%s err=%v",
			req.UserId,
			err,
		)
		return &pb.CreateProfileResponse{
			Success: false,
			Message: "failed to create profile",
		}, nil
	}

	return &pb.CreateProfileResponse{
		Success: true,
		Message: "profile created",
	}, nil
}

// ✅ 컴파일 타임 인터페이스 체크 (강력 추천)
var _ pb.UserServiceServer = (*UserGrpcServer)(nil)
