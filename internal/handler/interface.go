package handler

import "github.com/gin-gonic/gin"

// Handler는 모든 HTTP 요청 핸들러(Kakao, User, Post 등)가 반드시 구현해야 하는
// 공통 인터페이스입니다. 이를 통해 모든 핸들러를 HandlerMap에서 단일 타입으로 관리할 수 있습니다.
type Handler interface {
	// RegisterRoutes: 주어진 gin.Engine에 해당 핸들러의 모든 HTTP 엔드포인트를
	// 그룹핑하여 등록하는 역할을 수행합니다.
	RegisterRoutes(r *gin.Engine)
}
