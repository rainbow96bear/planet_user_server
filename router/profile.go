package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/user"
)

func RegisterProfileRoutes(r *gin.Engine) {

	profileGroup := r.Group("/profile")
	profileGroup.POST("/:userNickName", user.GetProfile)

	profileGroup.Use(middleware.AuthMiddleware()) // 미들웨어 적용
	{
		profileGroup.PATCH("/:userNickName", user.UpdateProfile)
	}
}
