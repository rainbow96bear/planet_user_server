package repository

import (
	"context"
	"time"

	"github.com/rainbow96bear/planet_utils/models"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

func (r *UsersRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	logger.Infof("UsersRepository: begin transaction")

	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Errorf("UsersRepository: failed to start tx: %v", tx.Error)
		return nil, tx.Error
	}

	logger.Infof("UsersRepository: transaction started successfully")
	return tx, nil
}

// 이메일로 유저 조회
func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*models.Users, error) {
	logger.Infof("UsersRepository:GetUserByEmail email=%s", email)

	var user models.Users
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		logger.Errorf("UsersRepository:GetUserByEmail failed email=%s error=%v", email, err)
		return nil, err
	}

	logger.Infof("UsersRepository:GetUserByEmail success user_id=%s", user.ID)
	return &user, nil
}

// UserID로 유저 조회
func (r *UsersRepository) GetUserByID(ctx context.Context, UserID string) (*models.Users, error) {
	logger.Infof("UsersRepository:GetUserByID user_id=%s", UserID)

	var user models.Users
	err := r.DB.WithContext(ctx).
		Where("id = ?", UserID).
		First(&user).Error

	if err != nil {
		logger.Errorf("UsersRepository:GetUserByID failed user_id=%s error=%v", UserID, err)
		return nil, err
	}

	logger.Infof("UsersRepository:GetUserByID success user_id=%s", user.ID)
	return &user, nil
}

// 유저 생성
func (r *UsersRepository) CreateUser(ctx context.Context, user *models.Users) error {
	logger.Infof("UsersRepository:CreateUser email=%s", user.Email)

	err := r.DB.WithContext(ctx).Create(user).Error

	if err != nil {
		logger.Errorf("UsersRepository:CreateUser failed email=%s error=%v", user.Email, err)
		return err
	}

	logger.Infof("UsersRepository:CreateUser success user_id=%s", user.ID)
	return nil
}

// 마지막 로그인 업데이트
func (r *UsersRepository) UpdateLastLogin(ctx context.Context, UserID string) error {
	logger.Infof("UsersRepository:UpdateLastLogin user_id=%s", UserID)

	now := time.Now()
	err := r.DB.WithContext(ctx).Model(&models.Users{}).
		Where("id = ?", UserID).
		Update("last_login", &now).Error

	if err != nil {
		logger.Errorf("UsersRepository:UpdateLastLogin failed user_id=%s error=%v", UserID, err)
		return err
	}

	logger.Infof("UsersRepository:UpdateLastLogin success user_id=%s", UserID)
	return nil
}

// Role 변경
func (r *UsersRepository) UpdateUserRole(ctx context.Context, UserID string, role string) error {
	logger.Infof("UsersRepository:UpdateUserRole user_id=%s role=%s", UserID, role)

	err := r.DB.WithContext(ctx).Model(&models.Users{}).
		Where("id = ?", UserID).
		Update("role", role).Error

	if err != nil {
		logger.Errorf("UsersRepository:UpdateUserRole failed user_id=%s role=%s error=%v", UserID, role, err)
		return err
	}

	logger.Infof("UsersRepository:UpdateUserRole success user_id=%s role=%s", UserID, role)
	return nil
}
