package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/user"
)

func RegisterProfileRoutes(r *gin.Engine, profileHandler *handler.ProfileHandler) {

	profileGroup := r.Group("/profile")
	profileGroup.GET("/:userNickName", profileHandler.GetProfileInfo)

	profileGroup.Use(middleware.AuthMiddleware(profileHandler.ProfileService)) // 미들웨어 적용
	{
		profileGroup.PATCH("/:userNickName", profileHandler.UpdateProfile)
	}
}
