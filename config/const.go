package config

var (
	PORT           string
	USER_GRPC_PORT string

	AUTH_GRPC_SERVER_ADDR      string
	ANALYTICS_GRPC_SERVER_ADDR string

	LOG_LEVEL           int16
	DB_GRPC_SERVER_ADDR string
	JWT_SECRET_KEY      string

	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string

	MaxTodoLength int16
)
