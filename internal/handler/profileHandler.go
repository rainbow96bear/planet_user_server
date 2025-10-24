package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
	AuthService    *service.AuthService
}

func NewProfileHandler(profileService *service.ProfileService, authService *service.AuthService) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
		AuthService:    authService,
	}
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	profileGroup := r.Group("/profile")
	profileGroup.GET("/:nickname", h.GetProfileInfo)
	profileGroup.Use(middleware.AuthMiddleware(h.AuthService))
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
		logger.Warnf("fail to get %s's profile info", nickname)
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

	authUuidValue, exists := c.Get("userUuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	authUuid, ok := authUuidValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid type"})
		return
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
