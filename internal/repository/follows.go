package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type FollowsRepository struct {
	DB *sql.DB
}

func (r *FollowsRepository) IsFollow(ctx context.Context, followerUuid, followeeUuid string) (bool, error) {
	logger.Infof("start to check %s follow %s", followerUuid, followeeUuid)
	defer logger.Infof("end to check %s follow %s", followerUuid, followeeUuid)

	query := `SELECT COUNT(*) FROM follows WHERE follower_uuid=? AND followee_uuid=?`

	var count int
	err := r.DB.QueryRowContext(ctx, query, followerUuid, followeeUuid).Scan(&count)

	return count > 0, err
}

func (r *FollowsRepository) FollowTx(ctx context.Context, tx *sql.Tx, followerUuid, followeeUuid string) error {
	logger.Infof("start %s follow %s", followerUuid, followeeUuid)
	defer logger.Infof("end %s follow %s", followerUuid, followeeUuid)

	query := `INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)`
	_, err := tx.ExecContext(ctx,
		query,
		followerUuid, followeeUuid,
	)
	if err != nil {
		return fmt.Errorf("failed to insert follow: %w", err)
	}

	return nil
}

func (r *FollowsRepository) UnfollowTx(ctx context.Context, tx *sql.Tx, followerUuid, followeeUuid string) error {
	logger.Infof("start %s unfollow %s", followerUuid, followeeUuid)
	defer logger.Infof("end %s unfollow %s", followerUuid, followeeUuid)

	query := `DELETE FROM follows WHERE follower_id = ? AND followee_id = ?`
	result, err := tx.ExecContext(ctx,
		query,
		followerUuid, followeeUuid,
	)
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no follow relation found to delete")
	}

	return nil
}
