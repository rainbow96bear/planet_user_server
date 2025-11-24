package userInit

import (
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes and returns a PostgreSQL database connection.
func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul",
		config.DB_HOST,
		config.DB_PORT,
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_NAME,
	)
	logger.Debugf("PostgreSQL DSN: %s", dsn)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open GORM DB: %w", err)
	}

	// sql.DB 가져오기 → 커넥션 풀 설정
	sqlDB, err := gormDB.DB()
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

	logger.Infof("✅ Successfully connected to PostgreSQL [%s:%s/%s]", config.DB_HOST, config.DB_PORT, config.DB_NAME)

	DB = gormDB
	return gormDB, nil
}
