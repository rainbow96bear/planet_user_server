package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_utils/models"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type ProfilesRepository struct {
	DB *gorm.DB
}

func (r *ProfilesRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	logger.Infof("starting transaction for ProfilesRepository")
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Errorf("failed to start transaction: %v", tx.Error)
		return nil, tx.Error
	}
	logger.Infof("transaction started successfully")
	return tx, nil
}

func (r *ProfilesRepository) GetUserIDByNickname(ctx context.Context, nickname string) (uuid.UUID, error) {
	logger.Infof("start to get user UUID by nickname: %s", nickname)

	var p models.Profiles
	err := r.DB.WithContext(ctx).
		Where("nickname = ?", nickname).
		First(&p).Error
	if err != nil {
		logger.Errorf("failed to get user UUID by nickname: %s", err.Error())
		return uuid.Nil, err
	}

	return p.UserID, nil
}

func (r *ProfilesRepository) GetFollowCounts(ctx context.Context, UserID uuid.UUID) (followerCount int, followingCount int, err error) {
	var profile models.Profiles
	if err := r.DB.WithContext(ctx).
		Select("follower_count", "following_count").
		Where("user_id = ?", UserID).
		First(&profile).Error; err != nil {
		return 0, 0, err
	}

	return profile.FollowerCount, profile.FollowingCount, nil
}

// user_id로 프로필 조회 (내 프로필)
func (r *ProfilesRepository) GetProfileByUserID(ctx context.Context, UserID string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get my profile info user_id: %s", UserID)

	var p models.Profiles
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", UserID).
		First(&p).Error
	if err != nil {
		logger.Errorf("failed to get my profile info: %s", err.Error())
		return nil, err
	}

	return &dto.ProfileInfo{
		UserID:       p.UserID,
		Nickname:     p.Nickname,
		Bio:          p.Bio,
		ProfileImage: p.ProfileImage,
		Theme:        p.Theme,
	}, nil
}

// 닉네임으로 프로필 조회
func (r *ProfilesRepository) GetProfileByNickname(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	logger.Infof("start to get profile info: %s", nickname)

	var p models.Profiles
	err := r.DB.WithContext(ctx).
		Where("nickname = ?", nickname).
		First(&p).Error
	if err != nil {
		logger.Errorf("failed to get profile info: %s", err.Error())
		return nil, err
	}

	return &dto.ProfileInfo{
		UserID:       p.UserID,
		Nickname:     p.Nickname,
		Bio:          p.Bio,
		ProfileImage: p.ProfileImage,
		Theme:        p.Theme,
	}, nil
}

// 프로필 업데이트
func (r *ProfilesRepository) UpdateProfile(ctx context.Context, profile *dto.ProfileUpdate) error {
	logger.Infof("update profile for user_id: %s", profile.UserID)

	return r.DB.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", profile.UserID).
		Updates(map[string]any{
			"nickname":      profile.Nickname,
			"bio":           profile.Bio,
			"profile_image": profile.ProfileImage,
		}).Error
}

// 테마 설정(JSONB)
func (r *ProfilesRepository) UpdateTheme(ctx context.Context, UserID string, theme map[string]interface{}) error {
	return r.DB.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", UserID).
		Update("theme", theme).Error
}

// 팔로워/팔로잉 증가(트랜잭션)
func (r *ProfilesRepository) IncrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerID, followeeID uuid.UUID) error {
	if err := tx.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", followerID).
		UpdateColumn("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", followeeID).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error; err != nil {
		return err
	}

	return nil
}

// 팔로워/팔로잉 감소(트랜잭션)
func (r *ProfilesRepository) DecrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerID, followeeID uuid.UUID) error {
	if err := tx.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", followerID).
		UpdateColumn("following_count", gorm.Expr("GREATEST(following_count - 1, 0)")).Error; err != nil {
		return err
	}

	if err := tx.WithContext(ctx).Model(&models.Profiles{}).
		Where("user_id = ?", followeeID).
		UpdateColumn("follower_count", gorm.Expr("GREATEST(follower_count - 1, 0)")).Error; err != nil {
		return err
	}

	return nil
}

// 테마 조회 (preset만 반환)
func (r *ProfilesRepository) GetTheme(ctx context.Context, userID uuid.UUID) (string, error) {
	logger.Infof("ProfilesRepository:GetTheme user_id=%s", userID)

	var p models.Profiles
	if err := r.DB.WithContext(ctx).
		Select("theme").
		Where("user_id = ?", userID).
		First(&p).Error; err != nil {
		logger.Errorf("ProfilesRepository:GetTheme failed user_id=%s error=%v", userID, err)
		return "light", err
	}

	// Theme JSON가 nil이거나 preset이 없는 경우 기본값으로 "light"
	theme := "light"
	if p.Theme != "" {
		theme = p.Theme
	}

	return theme, nil
}

// 테마 업데이트 (preset만 갱신)
func (r *ProfilesRepository) SetTheme(ctx context.Context, userID uuid.UUID, theme string) error {
	logger.Infof("ProfilesRepository:SetTheme user_id=%s theme=%v", userID, theme)

	if err := r.DB.WithContext(ctx).
		Model(&models.Profiles{}).
		Where("user_id = ?", userID).
		Update("theme", theme).Error; err != nil {
		logger.Errorf("ProfilesRepository:SetTheme failed user_id=%s theme=%s error=%v", userID, theme, err)
		return err
	}

	logger.Infof("ProfilesRepository:SetTheme success user_id=%s theme=%s", userID, theme)
	return nil
}

func (r *ProfilesRepository) IsMyProfile(ctx context.Context, UserID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.Profiles{}).
		Where("user_id = ?", UserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 닉네임으로 프로필 조회
func (r *ProfilesRepository) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
	logger.Infof("ProfilesRepository:GetProfileInfo nickname=%s", nickname)

	var p models.Profiles
	err := r.DB.WithContext(ctx).
		Where("nickname = ?", nickname).
		First(&p).Error
	if err != nil {
		logger.Errorf("ProfilesRepository:GetProfileInfo failed nickname=%s error=%v", nickname, err)
		return nil, err
	}

	info := &dto.ProfileInfo{
		UserID:         p.UserID,
		Nickname:       p.Nickname,
		Bio:            p.Bio,
		ProfileImage:   p.ProfileImage,
		Theme:          p.Theme,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
	}

	logger.Infof("ProfilesRepository:GetProfileInfo success nickname=%s", nickname)
	return info, nil
}

// UserID로 내 프로필 조회
func (r *ProfilesRepository) GetMyProfileInfo(ctx context.Context, UserID uuid.UUID) (*dto.ProfileInfo, error) {
	logger.Infof("ProfilesRepository:GetMyProfileInfo user_id=%s", UserID)

	var p models.Profiles
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", UserID).
		First(&p).Error
	if err != nil {
		logger.Errorf("ProfilesRepository:GetMyProfileInfo failed user_id=%s error=%v", UserID, err)
		return nil, err
	}

	info := &dto.ProfileInfo{
		UserID:         p.UserID,
		Nickname:       p.Nickname,
		Bio:            p.Bio,
		ProfileImage:   p.ProfileImage,
		Theme:          p.Theme,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
	}

	logger.Infof("ProfilesRepository:GetMyProfileInfo success user_id=%s", UserID)
	return info, nil
}
