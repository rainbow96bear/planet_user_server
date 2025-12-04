package router

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter는 Gin 엔진을 설정하고 반환하며, 라우트 등록을 위한 콜백 함수를 받습니다.
func SetupRouter(registerRoutes func(r *gin.Engine)) *gin.Engine {
	// 1. Gin 인스턴스 초기화
	r := gin.Default()

	// 2. 콜백 함수 실행 (여기서 GraphQL 라우트가 등록됨)
	registerRoutes(r)

	// 필요하다면 전역 미들웨어를 여기서 추가합니다.
	// r.Use(gin.Recovery())

	return r
}
