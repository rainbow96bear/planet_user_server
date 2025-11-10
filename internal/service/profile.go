package service

import (
	"context"
	"fmt"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type ProfileService struct {
	UsersRepo *repository.UsersRepository
}

func (s *ProfileService) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	profile, err := s.UsersRepo.GetProfileInfo(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("fail to get profile info")
	}

	userUuid, err := s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}

	followerCount, followingCount, err := s.UsersRepo.GetFollowCounts(ctx, userUuid)
	if err != nil {
		return nil, err
	}

	profile.FollowerCount = followerCount
	profile.FollowingCount = followingCount

	return profile, nil
}

func (s *ProfileService) GetMyProfileInfo(ctx context.Context, userUuid string) (*dto.ProfileInfo, error) {
	profile, err := s.UsersRepo.GetMyProfileInfo(ctx, userUuid)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("fail to get profile info")
	}

	return profile, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, profile *dto.ProfileInfo) error {
	err := s.UsersRepo.UpdateProfile(ctx, profile)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProfileService) GetUserUuidByNickname(ctx context.Context, nickname string) (string, error) {
	return s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
}

func (s *ProfileService) GetFollowCounts(ctx context.Context, userUuid string) (uint, uint, error) {
	followerCounts, followeeCounts, err := s.UsersRepo.GetFollowCounts(ctx, userUuid)
	if err != nil {
		return 0, 0, err
	}
	return followerCounts, followeeCounts, nil
}

func (s *ProfileService) IsMyProfile(ctx context.Context, userUUID string, nickname string) (bool, error) {
	profile, err := s.UsersRepo.GetMyProfileInfo(ctx, userUUID)
	if err != nil {
		return false, err
	}

	// 닉네임 일치 여부 확인
	return profile.Nickname == nickname, nil
}
