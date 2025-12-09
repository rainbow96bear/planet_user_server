package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

// NewTokensRepository는 TokensRepository 인스턴스를 생성합니다.
func NewProfilesRepository(db *gorm.DB) *ProfileRepository {
	if db == nil {
		panic("database connection is required")
	}
	return &ProfileRepository{
		db: db,
	}
}

// 헬퍼 함수: Context에서 트랜잭션을 확인하고, 있으면 Tx 객체를, 없으면 기본 DB 객체를 반환합니다.
func (r *ProfileRepository) getDB(ctx context.Context) *gorm.DB {
	// tx 패키지를 사용하여 Context에서 트랜잭션을 추출합니다.
	if tx := tx.GetTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx) // 기본 DB 연결 반환
}

// func (r *ProfileRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
// 	logger.Infof("starting transaction for ProfileRepository")
// 	tx := r.DB.WithContext(ctx).Begin()
// 	if tx.Error != nil {
// 		logger.Errorf("failed to start transaction: %v", tx.Error)
// 		return nil, tx.Error
// 	}
// 	logger.Infof("transaction started successfully")
// 	return tx, nil
// }

// func (r *ProfileRepository) GetUserIDByNickname(ctx context.Context, nickname string) (uuid.UUID, error) {
// 	logger.Infof("start to get user UUID by nickname: %s", nickname)

// 	var p models.Profile
// 	err := r.DB.WithContext(ctx).
// 		Where("nickname = ?", nickname).
// 		First(&p).Error
// 	if err != nil {
// 		logger.Errorf("failed to get user UUID by nickname: %s", err.Error())
// 		return uuid.Nil, err
// 	}

// 	return p.UserID, nil
// }

// func (r *ProfileRepository) GetFollowCounts(ctx context.Context, UserID uuid.UUID) (followerCount int, followingCount int, err error) {
// 	var profile models.Profile
// 	if err := r.DB.WithContext(ctx).
// 		Select("follower_count", "following_count").
// 		Where("user_id = ?", UserID).
// 		First(&profile).Error; err != nil {
// 		return 0, 0, err
// 	}

// 	return profile.FollowerCount, profile.FollowingCount, nil
// }

// // user_id로 프로필 조회 (내 프로필)
// func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, UserID string) (*dto.ProfileInfo, error) {
// 	logger.Infof("start to get my profile info user_id: %s", UserID)

// 	var p models.Profile
// 	err := r.DB.WithContext(ctx).
// 		Where("user_id = ?", UserID).
// 		First(&p).Error
// 	if err != nil {
// 		logger.Errorf("failed to get my profile info: %s", err.Error())
// 		return nil, err
// 	}

// 	return &dto.ProfileInfo{
// 		UserID:       p.UserID,
// 		Nickname:     p.Nickname,
// 		Bio:          p.Bio,
// 		ProfileImage: p.ProfileImage,
// 		Theme:        p.Theme,
// 	}, nil
// }

// // 닉네임으로 프로필 조회
// func (r *ProfileRepository) GetProfileByNickname(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
// 	logger.Infof("start to get profile info: %s", nickname)

// 	var p models.Profile
// 	err := r.DB.WithContext(ctx).
// 		Where("nickname = ?", nickname).
// 		First(&p).Error
// 	if err != nil {
// 		logger.Errorf("failed to get profile info: %s", err.Error())
// 		return nil, err
// 	}

// 	return &dto.ProfileInfo{
// 		UserID:       p.UserID,
// 		Nickname:     p.Nickname,
// 		Bio:          p.Bio,
// 		ProfileImage: p.ProfileImage,
// 		Theme:        p.Theme,
// 	}, nil
// }

// // 프로필 업데이트
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *dto.ProfileUpdate) error {
	db := r.getDB(ctx)
	logger.Infof("update profile for user_id: %s", profile.UserID)

	updates := utils.StructToUpdateMap(profile)

	if len(updates) == 0 {
		return nil
	}

	return db.WithContext(ctx).
		Model(&models.Profile{}).
		Where("user_id = ?", profile.UserID).
		Updates(updates).Error
}

// // 테마 설정(JSONB)
// func (r *ProfileRepository) UpdateTheme(ctx context.Context, UserID string, theme map[string]interface{}) error {
// 	return r.DB.WithContext(ctx).Model(&models.Profile{}).
// 		Where("user_id = ?", UserID).
// 		Update("theme", theme).Error
// }

// // 팔로워/팔로잉 증가(트랜잭션)
// func (r *ProfileRepository) IncrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerID, followeeID uuid.UUID) error {
// 	if err := tx.WithContext(ctx).Model(&models.Profile{}).
// 		Where("user_id = ?", followerID).
// 		UpdateColumn("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
// 		return err
// 	}

// 	if err := tx.WithContext(ctx).Model(&models.Profile{}).
// 		Where("user_id = ?", followeeID).
// 		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// // 팔로워/팔로잉 감소(트랜잭션)
// func (r *ProfileRepository) DecrementFollowCountsTx(ctx context.Context, tx *gorm.DB, followerID, followeeID uuid.UUID) error {
// 	if err := tx.WithContext(ctx).Model(&models.Profile{}).
// 		Where("user_id = ?", followerID).
// 		UpdateColumn("following_count", gorm.Expr("GREATEST(following_count - 1, 0)")).Error; err != nil {
// 		return err
// 	}

// 	if err := tx.WithContext(ctx).Model(&models.Profile{}).
// 		Where("user_id = ?", followeeID).
// 		UpdateColumn("follower_count", gorm.Expr("GREATEST(follower_count - 1, 0)")).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// // 테마 조회 (preset만 반환)
// func (r *ProfileRepository) GetTheme(ctx context.Context, userID uuid.UUID) (string, error) {
// 	logger.Infof("ProfileRepository:GetTheme user_id=%s", userID)

// 	var p models.Profile
// 	if err := r.DB.WithContext(ctx).
// 		Select("theme").
// 		Where("user_id = ?", userID).
// 		First(&p).Error; err != nil {
// 		logger.Errorf("ProfileRepository:GetTheme failed user_id=%s error=%v", userID, err)
// 		return "light", err
// 	}

// 	// Theme JSON가 nil이거나 preset이 없는 경우 기본값으로 "light"
// 	theme := "light"
// 	if p.Theme != "" {
// 		theme = p.Theme
// 	}

// 	return theme, nil
// }

// // 테마 업데이트 (preset만 갱신)
// func (r *ProfileRepository) SetTheme(ctx context.Context, userID uuid.UUID, theme string) error {
// 	logger.Infof("ProfileRepository:SetTheme user_id=%s theme=%v", userID, theme)

// 	if err := r.DB.WithContext(ctx).
// 		Model(&models.Profile{}).
// 		Where("user_id = ?", userID).
// 		Update("theme", theme).Error; err != nil {
// 		logger.Errorf("ProfileRepository:SetTheme failed user_id=%s theme=%s error=%v", userID, theme, err)
// 		return err
// 	}

// 	logger.Infof("ProfileRepository:SetTheme success user_id=%s theme=%s", userID, theme)
// 	return nil
// }

func (r *ProfileRepository) IsMyProfile(ctx context.Context, UserID uuid.UUID) (bool, error) {
	db := r.getDB(ctx)
	var count int64
	err := db.WithContext(ctx).
		Model(&models.Profile{}).
		Where("user_id = ?", UserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// // 닉네임으로 프로필 조회
// func (r *ProfileRepository) GetProfileInfo(ctx context.Context, nickname string) (*dto.ProfileInfo, error) {
// 	logger.Infof("ProfileRepository:GetProfileInfo nickname=%s", nickname)

// 	var p models.Profile
// 	err := r.DB.WithContext(ctx).
// 		Where("nickname = ?", nickname).
// 		First(&p).Error
// 	if err != nil {
// 		logger.Errorf("ProfileRepository:GetProfileInfo failed nickname=%s error=%v", nickname, err)
// 		return nil, err
// 	}

// 	info := &dto.ProfileInfo{
// 		UserID:         p.UserID,
// 		Nickname:       p.Nickname,
// 		Bio:            p.Bio,
// 		ProfileImage:   p.ProfileImage,
// 		Theme:          p.Theme,
// 		FollowerCount:  p.FollowerCount,
// 		FollowingCount: p.FollowingCount,
// 	}

// 	logger.Infof("ProfileRepository:GetProfileInfo success nickname=%s", nickname)
// 	return info, nil
// }

// // UserID로 내 프로필 조회
func (r *ProfileRepository) GetMyProfileInfo(ctx context.Context, userID uuid.UUID) (*models.Profile, error) {
	db := r.getDB(ctx)
	var p models.Profile
	if err := db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
