package service

import (
	"context"

	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type FollowService struct {
	UsersRepo   *repository.UsersRepository
	FollowsRepo *repository.FollowsRepository
}

func (s *FollowService) IsFollow(ctx context.Context, followerUuid, followeeUuid string) (bool, error) {
	isFollow, err := s.FollowsRepo.IsFollow(ctx, followerUuid, followeeUuid)
	if err != nil {
		return false, err
	}

	return isFollow, nil
}

func (s *FollowService) Follow(ctx context.Context, followerUuid, followeeUuid string) error {
	// GORM 트랜잭션 시작
	tx := s.UsersRepo.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 실패하면 rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 팔로우 생성
	if err := s.FollowsRepo.FollowTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		tx.Rollback()
		return err
	}

	// 팔로우 카운트 증가
	if err := s.UsersRepo.IncrementFollowCountsTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *FollowService) Unfollow(ctx context.Context, followerUuid, followeeUuid string) error {
	tx := s.UsersRepo.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.FollowsRepo.UnfollowTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.UsersRepo.DecrementFollowCountsTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
