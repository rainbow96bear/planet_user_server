package repository

import (
	"context"
	"fmt"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_utils/model"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

// 닉네임으로 user_uuid 조회
func (r *UsersRepository) GetUserUuidByNickname(ctx context.Context, nickname string) (string, error) {
	logger.Infof("start to get user uuid: %s", nickname)
	defer logger.Infof("end to get user uuid: %s", nickname)

	var user model.User
	err := r.DB.WithContext(ctx).Select("user_uuid").Where("nickname = ?", nickname).First(&user).Error
	if err != nil {
		logger.Errorf("failed to get user uuid ERR[%s]", err.Error())
		return "", err
	}

	return user.UserUUID, nil
}

// 닉네임으로 프로필 조회
func (r *UsersRepository) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get profile info : %s", nickname)
	defer logger.Infof("end to get profile info: %s", nickname)

	var user model.User
	err := r.DB.WithContext(ctx).
		Select("user_uuid", "nickname", "profile_image", "bio", "email").
		Where("nickname = ?", nickname).
		First(&user).Error
	if err != nil {
		logger.Errorf("failed to get profile info ERR[%s]", err.Error())
		return nil, err
	}

	return &dto.ProfileInfo{
		UserUUID:     user.UserUUID,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		Email:        user.Email,
	}, nil
}

// user_uuid로 내 프로필 조회
func (r *UsersRepository) GetMyProfileInfo(ctx context.Context, userUUID string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get profile info user uuid : %s", userUUID)
	defer logger.Infof("end to get profile info user uuid : %s", userUUID)

	var user model.User
	err := r.DB.WithContext(ctx).
		Select("user_uuid", "nickname", "profile_image", "bio", "email").
		Where("user_uuid = ?", userUUID).
		First(&user).Error
	if err != nil {
		logger.Errorf("failed to get profile info ERR[%s]", err.Error())
		return nil, err
	}

	return &dto.ProfileInfo{
		UserUUID:     user.UserUUID,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		Email:        user.Email,
	}, nil
}

// 프로필 업데이트
func (r *UsersRepository) UpdateProfile(ctx context.Context, profile *dto.ProfileInfo) error {
	logger.Infof("start to update profile info: %s", profile.Nickname)
	defer logger.Infof("end to update profile info: %s", profile.Nickname)

	err := r.DB.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", profile.UserUUID).
		Updates(map[string]any{
			"nickname":      profile.Nickname,
			"profile_image": profile.ProfileImage,
			"bio":           profile.Bio,
			"email":         profile.Email,
		}).Error
	if err != nil {
		logger.Errorf("failed to update profile info ERR[%s]", err.Error())
		return err
	}

	return nil
}

// 팔로우 수 증가
func (r *UsersRepository) IncrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerUuid, followeeUuid string) error {
	if err := tx.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", followerUuid).
		UpdateColumn("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment following_count: %w", err)
	}

	if err := tx.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", followeeUuid).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment follower_count: %w", err)
	}

	return nil
}

// 팔로우 수 감소
func (r *UsersRepository) DecrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerUuid, followeeUuid string) error {
	if err := tx.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", followerUuid).
		UpdateColumn("following_count", gorm.Expr("CASE WHEN following_count > 0 THEN following_count - 1 ELSE 0 END")).Error; err != nil {
		return fmt.Errorf("failed to decrement following_count: %w", err)
	}

	if err := tx.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", followeeUuid).
		UpdateColumn("follower_count", gorm.Expr("CASE WHEN follower_count > 0 THEN follower_count - 1 ELSE 0 END")).Error; err != nil {
		return fmt.Errorf("failed to decrement follower_count: %w", err)
	}

	return nil
}

// 팔로워/팔로잉 수 조회
func (r *UsersRepository) GetFollowCounts(ctx context.Context, userUuid string) (uint, uint, error) {
	logger.Infof("start to get follow counts: %s", userUuid)
	defer logger.Infof("end to  get follow counts: %s", userUuid)

	var user model.User
	err := r.DB.WithContext(ctx).Select("follower_count", "following_count").Where("user_uuid = ?", userUuid).First(&user).Error
	if err != nil {
		return 0, 0, err
	}

	return user.FollowerCount, user.FollowingCount, nil
}

// 테마 조회
func (r *UsersRepository) GetTheme(ctx context.Context, userUuid string) (string, error) {
	logger.Infof("start to get theme user_uuid : %s", userUuid)
	defer logger.Infof("end to get theme user_uuid : %s", userUuid)

	var user model.User
	err := r.DB.WithContext(ctx).Select("theme").Where("user_uuid = ?", userUuid).First(&user).Error
	if err != nil {
		return "", err
	}

	return user.Theme, nil
}

// 테마 설정
func (r *UsersRepository) SetTheme(ctx context.Context, userUuid, theme string) error {
	logger.Infof("start to set theme : %s, user_uuid : %s", theme, userUuid)
	defer logger.Infof("end to set theme : %s, user_uuid : %s", theme, userUuid)

	err := r.DB.WithContext(ctx).Model(&model.User{}).
		Where("user_uuid = ?", userUuid).
		Update("theme", theme).Error
	if err != nil {
		logger.Errorf("failed to update theme for user_uuid %s: %v", userUuid, err)
		return err
	}

	return nil
}
