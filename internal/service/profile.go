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
		return nil, fmt.Errorf("fail to get profile info ERR[%s]", err.Error())
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
