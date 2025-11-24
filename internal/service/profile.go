package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type ProfileService struct {
	ProfilesRepo *repository.ProfilesRepository
}

// 닉네임으로 사용자 UUID 조회
func (s *ProfileService) GetUserIDByNickname(ctx context.Context, nickname string) (uuid.UUID, error) {
	UserID, err := s.ProfilesRepo.GetUserIDByNickname(ctx, nickname)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user ID by nickname: %w", err)
	}
	if UserID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("user not found for nickname: %s", nickname)
	}
	return UserID, nil
}

// 다른 유저 프로필 조회
func (s *ProfileService) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	profile, err := s.ProfilesRepo.GetProfileInfo(ctx, nickname)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile info: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("profile not found for nickname: %s", nickname)
	}
	return profile, nil
}

func (s *ProfileService) GetFollowCounts(ctx context.Context, UserID uuid.UUID) (followerCount int, followingCount int, err error) {
	return s.ProfilesRepo.GetFollowCounts(ctx, UserID)
}

// 내 프로필 조회
func (s *ProfileService) GetMyProfileInfo(ctx context.Context, UserID uuid.UUID) (*dto.ProfileInfo, error) {
	profile, err := s.ProfilesRepo.GetMyProfileInfo(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my profile info: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("my profile not found for user: %s", UserID)
	}
	return profile, nil
}

// 프로필 업데이트
func (s *ProfileService) UpdateProfile(ctx context.Context, UserID uuid.UUID, nickname string, req *dto.ProfileUpdateRequest) (*dto.ProfileInfo, error) {
	// 먼저 UUID와 닉네임 일치 여부 검증
	isMyProfile, err := s.ProfilesRepo.IsMyProfile(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify profile ownership: %w", err)
	}
	if !isMyProfile {
		return nil, fmt.Errorf("unauthorized: cannot update another user's profile")
	}

	// DTO -> 내부 모델 변환
	updateModel := dto.ToProfileUpdateModel(req, UserID)

	// 업데이트
	if err := s.ProfilesRepo.UpdateProfile(ctx, updateModel); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// 업데이트 후 최신 프로필 반환
	return s.GetMyProfileInfo(ctx, UserID)
}

// 테마 조회
func (s *ProfileService) GetTheme(ctx context.Context, UserID uuid.UUID) (string, error) {
	theme, err := s.ProfilesRepo.GetTheme(ctx, UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get theme: %w", err)
	}
	return theme, nil
}

// 테마 설정
func (s *ProfileService) SetTheme(ctx context.Context, userID uuid.UUID, theme string) error {
	if err := s.ProfilesRepo.SetTheme(ctx, userID, theme); err != nil {
		return fmt.Errorf("failed to set theme: %w", err)
	}
	return nil
}
