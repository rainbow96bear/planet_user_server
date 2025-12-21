package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_utils/models"
	"gorm.io/gorm"
)

type FollowsRepository struct {
	db *gorm.DB
}

func NewFollowsRepository(db *gorm.DB) *FollowsRepository {
	if db == nil {
		panic("database connection is required")
	}
	return &FollowsRepository{db: db}
}

func (r *FollowsRepository) getDB(ctx context.Context) *gorm.DB {
	if txDB := tx.GetTx(ctx); txDB != nil {
		return txDB.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

//
// =======================
// Query
// =======================
//

// follower → followee 팔로우 여부
func (r *FollowsRepository) IsFollowing(
	ctx context.Context,
	followerID uuid.UUID,
	followeeID uuid.UUID,
) (bool, error) {

	db := r.getDB(ctx)

	var count int64
	if err := db.
		Model(&models.Follows{}).
		Where("follower_uuid = ? AND followee_uuid = ?", followerID, followeeID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

//
// =======================
// Command
// =======================
//

// 팔로우 생성
func (r *FollowsRepository) Create(
	ctx context.Context,
	followerID uuid.UUID,
	followeeID uuid.UUID,
) error {

	db := r.getDB(ctx)

	follow := &models.Follows{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	if err := db.Create(follow).Error; err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	return nil
}

// 팔로우 삭제
func (r *FollowsRepository) Delete(
	ctx context.Context,
	followerID uuid.UUID,
	followeeID uuid.UUID,
) error {

	db := r.getDB(ctx)

	result := db.
		Where("follower_uuid = ? AND followee_uuid = ?", followerID, followeeID).
		Delete(&models.Follows{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete follow: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("follow relation not found")
	}

	return nil
}
