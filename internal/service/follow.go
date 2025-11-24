package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type FollowService struct {
	ProfilesRepo *repository.ProfilesRepository
	FollowsRepo  *repository.FollowsRepository
}

func (s *FollowService) IsFollow(ctx context.Context, followerUuid, followeeUuid uuid.UUID) (bool, error) {
	isFollow, err := s.FollowsRepo.IsFollow(ctx, followerUuid, followeeUuid)
	if err != nil {
		return false, err
	}

	return isFollow, nil
}

func (s *FollowService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// GORM 트랜잭션 시작

	tx, err := s.ProfilesRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// 트랜잭션 rollback 보장
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 팔로우 생성
	if err := s.FollowsRepo.FollowTx(ctx, tx, followerID, followingID); err != nil {
		tx.Rollback()
		return err
	}

	// 팔로우 카운트 증가
	if err := s.ProfilesRepo.IncrementFollowCountsTx(ctx, tx, followerID, followingID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *FollowService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	tx, err := s.ProfilesRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// 트랜잭션 rollback 보장
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := s.FollowsRepo.UnfollowTx(ctx, tx, followerID, followingID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.ProfilesRepo.DecrementFollowCountsTx(ctx, tx, followerID, followingID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
