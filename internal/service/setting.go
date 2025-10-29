package service

import (
	"context"

	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type SettingService struct {
	UsersRepo *repository.UsersRepository
}

func (s *SettingService) GetTheme(ctx context.Context, userUuid string) (string, error) {
	theme, err := s.UsersRepo.GetTheme(ctx, userUuid)
	if err != nil {
		return "", err
	}

	return theme, nil
}

func (s *SettingService) SetTheme(ctx context.Context, userUuid, theme string) error {
	err := s.UsersRepo.SetTheme(ctx, userUuid, theme)
	if err != nil {
		return err
	}

	return nil
}
