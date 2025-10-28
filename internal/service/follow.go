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
	tx, err := s.UsersRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.FollowsRepo.FollowTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		return err
	}

	if err := s.UsersRepo.IncrementFollowCountsTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *FollowService) Unfollow(ctx context.Context, followerUuid, followeeUuid string) error {
	tx, err := s.UsersRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.FollowsRepo.UnfollowTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		return err
	}

	if err := s.UsersRepo.DecrementFollowCountsTx(ctx, tx, followerUuid, followeeUuid); err != nil {
		return err
	}

	return tx.Commit()
}
