package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type UsersRepository struct {
	DB *sql.DB
}

func (r *UsersRepository) GetUserUuidByNickname(ctx context.Context, nickname string) (string, error) {
	logger.Infof("start to get user uuid: %s", nickname)
	defer logger.Infof("end to get user uuid: %s", nickname)

	query := `SELECT user_uuid FROM users WHERE nickname = ?`

	var userUuid string
	err := r.DB.QueryRowContext(ctx, query, nickname).Scan(&userUuid)
	if err != nil {
		logger.Errorf("failed to get user uuid ERR[%s]", err.Error())
		return "", err
	}

	return userUuid, nil
}

func (r *UsersRepository) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get profile info : %s", nickname)
	defer logger.Infof("end to get profile info: %s", nickname)

	query := `SELECT user_uuid, nickname, profile_image, bio, email FROM users WHERE nickname = ?`

	profileInfo := &dto.ProfileInfo{}
	err := r.DB.QueryRowContext(ctx, query, nickname).Scan(
		&profileInfo.UserUuid,
		&profileInfo.Nickname,
		&profileInfo.ProfileImage,
		&profileInfo.Bio,
		&profileInfo.Email,
	)
	if err != nil {
		logger.Errorf("failed to get profile info ERR[%s]", err.Error())
		return nil, err
	}

	return profileInfo, nil
}

func (r *UsersRepository) GetMyProfileInfo(ctx context.Context, userUuid string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get profile info user uuid : %s", userUuid)
	defer logger.Infof("end to get profile info user uuid : %s", userUuid)

	query := `SELECT user_uuid, nickname, profile_image, bio, email FROM users WHERE user_uuid = ?`

	profileInfo := &dto.ProfileInfo{}
	err := r.DB.QueryRowContext(ctx, query, userUuid).Scan(
		&profileInfo.UserUuid,
		&profileInfo.Nickname,
		&profileInfo.ProfileImage,
		&profileInfo.Bio,
		&profileInfo.Email,
	)
	if err != nil {
		logger.Errorf("failed to get profile info ERR[%s]", err.Error())
		return nil, err
	}

	return profileInfo, nil
}

func (r *UsersRepository) UpdateProfile(ctx context.Context, profile *dto.ProfileInfo) error {
	logger.Infof("start to update profile info: %s", profile.Nickname)
	defer logger.Infof("end to update profile info: %s", profile.Nickname)

	query := `
		UPDATE users
		SET nickname = ?, profile_image = ?, bio = ?, email = ?
		WHERE user_uuid = ?
	`

	_, err := r.DB.ExecContext(ctx, query,
		profile.Nickname,
		profile.ProfileImage,
		profile.Bio,
		profile.Email,
		profile.UserUuid,
	)
	if err != nil {
		logger.Errorf("failed to update profile info ERR[%s]", err.Error())
		return err
	}

	return nil
}

func (r *UsersRepository) IncrementFollowCountsTx(ctx context.Context, tx *sql.Tx, followerUuid, followeeUuid string) error {
	_, err := tx.ExecContext(ctx,
		"UPDATE users SET following_count = following_count + 1 WHERE uuid = ?",
		followerUuid,
	)
	if err != nil {
		return fmt.Errorf("failed to increment following_count: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE users SET follower_count = follower_count + 1 WHERE uuid = ?",
		followeeUuid,
	)
	if err != nil {
		return fmt.Errorf("failed to increment follower_count: %w", err)
	}

	return nil
}

func (r *UsersRepository) DecrementFollowCountsTx(ctx context.Context, tx *sql.Tx, followerUuid, followeeUuid string) error {
	// following_count 감소 (0 미만 방지)
	followerQuery := `
		UPDATE users 
		SET following_count = CASE WHEN following_count > 0 THEN following_count - 1 ELSE 0 END 
		WHERE uuid = ?
		`
	_, err := tx.ExecContext(ctx, followerQuery, followerUuid)
	if err != nil {
		return fmt.Errorf("failed to decrement following_count: %w", err)
	}

	// follower_count 감소 (0 미만 방지)
	followeeQuery := `
		UPDATE users 
		SET follower_count = CASE WHEN follower_count > 0 THEN follower_count - 1 ELSE 0 END 
		WHERE uuid = ?
		`
	_, err = tx.ExecContext(ctx, followeeQuery, followeeUuid)
	if err != nil {
		return fmt.Errorf("failed to decrement follower_count: %w", err)
	}

	return nil
}

func (r *UsersRepository) GetFollowCounts(ctx context.Context, userUuid string) (uint, uint, error) {
	logger.Infof("start to get follow counts: %s", userUuid)
	defer logger.Infof("end to  get follow counts: %s", userUuid)

	var followCount, followingCount uint
	query := `SELECT follower_count, following_count FROM users WHERE user_uuid = ?`
	err := r.DB.QueryRowContext(ctx, query, userUuid).Scan(&followCount, &followingCount)

	if err != nil {
		return 0, 0, err
	}

	return followCount, followingCount, nil
}

func (r *UsersRepository) GetTheme(ctx context.Context, userUuid string) (string, error) {
	logger.Infof("start to get theme user_uuid : %s", userUuid)
	defer logger.Infof("end to get theme user_uuid : %s", userUuid)

	var theme string
	query := `SELECT theme FROM users WHERE user_uuid = ?`
	err := r.DB.QueryRowContext(ctx, query, userUuid).Scan(&theme)

	if err != nil {
		return "", err
	}

	return theme, nil
}

func (r *UsersRepository) SetTheme(ctx context.Context, userUuid, theme string) error {
	logger.Infof("start to set theme : %s, user_uuid : %s", theme, userUuid)
	defer logger.Infof("end to set theme : %s, user_uuid : %s", theme, userUuid)

	query := `UPDATE users SET theme = ? WHERE user_uuid = ?`
	_, err := r.DB.ExecContext(ctx, query, theme, userUuid)
	if err != nil {
		logger.Errorf("failed to update theme for user_uuid %s: %v", userUuid, err)
		return err
	}

	return nil
}
