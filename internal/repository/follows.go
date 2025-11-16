package repository

import (
	"context"
	"fmt"

	"github.com/rainbow96bear/planet_utils/model"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type FollowsRepository struct {
	DB *gorm.DB
}

// 팔로우 관계 존재 여부 확인
func (r *FollowsRepository) IsFollow(ctx context.Context, followerUUID, followeeUUID string) (bool, error) {
	logger.Infof("start to check %s follow %s", followerUUID, followeeUUID)
	defer logger.Infof("end to check %s follow %s", followerUUID, followeeUUID)

	var count int64
	err := r.DB.WithContext(ctx).
		Model(&model.Follows{}).
		Where("follower_uuid = ? AND followee_uuid = ?", followerUUID, followeeUUID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 팔로우 생성 (트랜잭션 지원)
func (r *FollowsRepository) FollowTx(ctx context.Context, tx *gorm.DB, followerUUID, followeeUUID string) error {
	logger.Infof("start %s follow %s", followerUUID, followeeUUID)
	defer logger.Infof("end %s follow %s", followerUUID, followeeUUID)

	follow := &model.Follows{
		FollowerUUID: followerUUID,
		FolloweeUUID: followeeUUID,
	}

	if err := tx.WithContext(ctx).Create(follow).Error; err != nil {
		return fmt.Errorf("failed to insert follow: %w", err)
	}

	return nil
}

// 팔로우 삭제 (트랜잭션 지원)
func (r *FollowsRepository) UnfollowTx(ctx context.Context, tx *gorm.DB, followerUUID, followeeUUID string) error {
	logger.Infof("start %s unfollow %s", followerUUID, followeeUUID)
	defer logger.Infof("end %s unfollow %s", followerUUID, followeeUUID)

	result := tx.WithContext(ctx).
		Where("follower_uuid = ? AND followee_uuid = ?", followerUUID, followeeUUID).
		Delete(&model.Follows{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete follow: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no follow relation found to delete")
	}

	return nil
}
