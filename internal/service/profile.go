package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/planet_err"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type ProfileServiceInterface interface {
	IsNicknameAvailable(ctx context.Context, nickname string) (bool, error)
	CreateProfile(ctx context.Context, req dto.CreateProfileRequest) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.ProfileUpdate) (*dto.UserProfile, error)
	GetMyProfileInfo(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error)
	GetUserProfileInfo(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error)
}

type ProfileService struct {
	DB           *gorm.DB
	ProfilesRepo *repository.ProfileRepository
}

func NewProfileService(
	db *gorm.DB,
	profilesRepo *repository.ProfileRepository,
) ProfileServiceInterface {
	return &ProfileService{
		DB:           db,
		ProfilesRepo: profilesRepo,
	}
}

// 닉네임 중복 검사
func (s *ProfileService) IsNicknameAvailable(ctx context.Context, nickname string) (bool, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Profile{}).Where("nickname = ?", nickname).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// 회원가입 시 Profile 생성 (트랜잭션 + DTO)
func (s *ProfileService) CreateProfile(ctx context.Context, req dto.CreateProfileRequest) (*dto.ProfileResponse, error) {
	logger.Debugf("ProfileService: CreateProfile attempt for nickname=%s", req.Nickname)

	// 0. 트랜잭션 시작
	txDB, newCtx, err := tx.BeginTx(ctx, s.DB)
	if err != nil {
		logger.Errorf("CreateProfile: failed to start transaction: %v", err)
		return nil, errors.New("failed to start transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("CreateProfile: panic occurred, rollback triggered: %v", r)
			txDB.Rollback()
			panic(r)
		}
	}()

	ctx = newCtx

	// 1. 닉네임 중복 체크
	available, err := s.IsNicknameAvailable(ctx, req.Nickname)
	if err != nil {
		txDB.Rollback()
		return nil, err
	}
	if !available {
		txDB.Rollback()
		return nil, errors.New("nickname already in use")
	}

	// 2. Profile 생성
	profile := &models.Profile{
		UserID:       req.UserID,
		Nickname:     req.Nickname,
		Bio:          req.Bio,
		ProfileImage: req.ProfileImage,
		Theme:        req.Theme, // 문자열 그대로 저장
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := txDB.Create(profile).Error; err != nil {
		txDB.Rollback()
		return nil, err
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		return nil, errors.New("failed to commit transaction")
	}

	resp := &dto.ProfileResponse{
		UserID:       profile.UserID,
		Nickname:     profile.Nickname,
		Bio:          profile.Bio,
		ProfileImage: profile.ProfileImage,
		Theme:        profile.Theme, // 이제 문자열 그대로 전달
		CreatedAt:    profile.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    profile.UpdatedAt.Format(time.RFC3339),
	}

	return resp, nil
}

// // 닉네임으로 사용자 UUID 조회
// func (s *ProfileService) GetUserIDByNickname(ctx context.Context, nickname string) (uuid.UUID, error) {
// 	UserID, err := s.ProfilesRepo.GetUserIDByNickname(ctx, nickname)
// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("failed to get user ID by nickname: %w", err)
// 	}
// 	if UserID == uuid.Nil {
// 		return uuid.Nil, fmt.Errorf("user not found for nickname: %s", nickname)
// 	}
// 	return UserID, nil
// }

// // 다른 유저 프로필 조회
// func (s *ProfileService) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
// 	profile, err := s.ProfilesRepo.GetProfileInfo(ctx, nickname)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get profile info: %w", err)
// 	}
// 	if profile == nil {
// 		return nil, fmt.Errorf("profile not found for nickname: %s", nickname)
// 	}
// 	return profile, nil
// }

// func (s *ProfileService) GetFollowCounts(ctx context.Context, UserID uuid.UUID) (followerCount int, followingCount int, err error) {
// 	return s.ProfilesRepo.GetFollowCounts(ctx, UserID)
// }

// // 내 프로필 조회
func (s *ProfileService) GetMyProfileInfo(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error) {
	profile, err := s.ProfilesRepo.GetMyProfileInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// 모델 → DTO 변환
	return &dto.UserProfile{
		ID:             profile.ID,
		UserID:         profile.UserID,
		Nickname:       profile.Nickname,
		Bio:            profile.Bio,
		ProfileImage:   profile.ProfileImage,
		FollowerCount:  profile.FollowerCount,
		FollowingCount: profile.FollowingCount,
		Theme:          profile.Theme,
	}, nil
}

func (s *ProfileService) GetUserProfileInfo(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error) {
	profile, err := s.ProfilesRepo.GetUserProfileInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// 모델 → DTO 변환
	return &dto.UserProfile{
		ID:             profile.ID,
		UserID:         profile.UserID,
		Nickname:       profile.Nickname,
		Bio:            profile.Bio,
		ProfileImage:   profile.ProfileImage,
		FollowerCount:  profile.FollowerCount,
		FollowingCount: profile.FollowingCount,
		Theme:          profile.Theme,
	}, nil
}

// // 프로필 업데이트
func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	req *dto.ProfileUpdate,
) (*dto.UserProfile, error) {

	// 소유권 확인
	isMyProfile, err := s.ProfilesRepo.IsMyProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify profile ownership: %w", err)
	}
	if !isMyProfile {
		return nil, fmt.Errorf("unauthorized: cannot update another user's profile")
	}

	// 업데이트
	if err := s.ProfilesRepo.UpdateProfile(ctx, req); err != nil {
		// 닉네임 중복 오류 처리
		if errors.Is(err, planet_err.ErrNicknameDuplicate) {
			return nil, planet_err.ErrNicknameDuplicate
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// 최신 프로필 반환
	profile, err := s.GetMyProfileInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated profile: %w", err)
	}

	return profile, nil
}

// // 테마 조회
// func (s *ProfileService) GetTheme(ctx context.Context, UserID uuid.UUID) (string, error) {
// 	theme, err := s.ProfilesRepo.GetTheme(ctx, UserID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get theme: %w", err)
// 	}
// 	return theme, nil
// }

// // 테마 설정
// func (s *ProfileService) SetTheme(ctx context.Context, userID uuid.UUID, theme string) error {
// 	if err := s.ProfilesRepo.SetTheme(ctx, userID, theme); err != nil {
// 		return fmt.Errorf("failed to set theme: %w", err)
// 	}
// 	return nil
// }
