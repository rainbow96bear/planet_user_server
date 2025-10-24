package repository

import (
	"context"
	"database/sql"

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

	query := `SELECT nickname, profile_image, bio, email FROM users WHERE nickname = ?`

	profileInfo := &dto.ProfileInfo{}
	err := r.DB.QueryRowContext(ctx, query, nickname).Scan(profileInfo)
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
