package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

func InitConfig(mode string) {
	var err error

	switch mode {
	case "prod":
		err = godotenv.Load("./env/.env.prod")
	case "dev":
		err = godotenv.Load("./env/.env.dev")
	}

	if err != nil {
		fmt.Println("[CONFIG] fail to load .env file, 기본값 dev 사용")
	}

	// default config
	PORT = getString("PORT")
	USER_GRPC_PORT = getString("USER_GRPC_PORT")

	AUTH_GRPC_SERVER_ADDR = getString("AUTH_GRPC_SERVER_ADDR")

	LOG_LEVEL = getInt16("LOG_LEVEL")
	DB_GRPC_SERVER_ADDR = getString("DB_GRPC_SERVER_ADDR")
	JWT_SECRET_KEY = getString("JWT_SECRET_KEY")

	DB_USER = getString("DB_USER")
	DB_PASSWORD = getString("DB_PASSWORD")
	DB_HOST = getString("DB_HOST")
	DB_PORT = getString("DB_PORT")
	DB_NAME = getString("DB_NAME")

	MaxTodoLength = getInt16("MaxTodoLength")
}

func getString(envName string) string {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}
	return v
}

func getInt16(envName string) int16 {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		logger.Errorf("[CONFIG] %s must be int, got %s\n", envName, v)
		os.Exit(1)
	}
	return int16(num)
}
