package userInit

import (
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// DSN 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)
	logger.Debugf(dsn)

	// GORM DB 연결
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 기본 SQL DB 객체 가져와 커넥션 풀 설정
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 연결 확인
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Infof("✅ Successfully connected to database [%s:%s/%s]", config.DB_HOST, config.DB_PORT, config.DB_NAME)

	DB = db
	return db, nil
}
