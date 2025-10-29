package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
	FollowService  *service.FollowService
}

func NewProfileHandler(profileService *service.ProfileService, followService *service.FollowService) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
		FollowService:  followService,
	}
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	profileGroup := r.Group("/profile")
	profileGroup.GET("/:nickname", h.GetProfileInfo)
	profileGroup.Use(middleware.AuthMiddleware())
	{
		profileGroup.PATCH("/:nickname", h.UpdateProfile)
	}
}

func (h *ProfileHandler) GetProfileInfo(c *gin.Context) {
	logger.Infof("start to get profile")
	defer logger.Infof("end to get profile")
	ctx := c.Request.Context()

	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	profileInfo, err := h.ProfileService.GetProfileInfo(ctx, nickname)
	if err != nil {
		logger.Warnf("fail to get %s's profile info ERR[%s]", nickname, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	profileResponse := *dto.ToProfileResponse(profileInfo)
	c.JSON(http.StatusOK, gin.H{
		"profile": profileResponse,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	logger.Infof("start to update profile")
	defer logger.Infof("end to update profile")

	ctx := c.Request.Context()

	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
	}

	profileUpdateRequest := &dto.ProfileUpdateRequest{}

	if err := c.ShouldBindJSON(profileUpdateRequest); err != nil {
		logger.Warnf("fail to bind to json about profileInfo : %+v", profileUpdateRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profileInfo := dto.ToProfileInfo(profileUpdateRequest, authUuid)

	if err := h.ProfileService.UpdateProfile(ctx, profileInfo); err != nil {
		logger.Warnf("fail to update profileInfo : %+v", profileInfo)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profileInfo,
	})
}
