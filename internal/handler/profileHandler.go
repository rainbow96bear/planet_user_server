package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		profileGroup.GET("/me", h.GetMyProfileInfo)
		profileGroup.PATCH("/:nickname", h.UpdateProfile)
		profileGroup.GET("/theme", h.GetTheme)
		profileGroup.POST("/theme", h.SetTheme)
	}
}

// ---------------------- Handler ----------------------

// 다른 유저 프로필 조회
func (h *ProfileHandler) GetProfileInfo(c *gin.Context) {
	logger.Infof("GetProfileInfo start")
	defer logger.Infof("GetProfileInfo end")

	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	profileInfo, err := h.ProfileService.GetProfileInfo(ctx, nickname)
	if err != nil {
		logger.Warnf("Failed to get profile info for %s: %v", nickname, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile info"})
		return
	}

	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
}

// 내 프로필 조회
func (h *ProfileHandler) GetMyProfileInfo(c *gin.Context) {
	logger.Infof("GetMyProfileInfo start")
	defer logger.Infof("GetMyProfileInfo end")

	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	if err != nil || authID == uuid.Nil {
		logger.Errorf("failed to get authenticated user UUID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profileInfo, err := h.ProfileService.GetMyProfileInfo(ctx, authID)
	if err != nil {
		logger.Warnf("Failed to get profile info for %s: %v", authID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile info"})
		return
	}

	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
}

// 프로필 업데이트
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	logger.Infof("UpdateProfile start")
	defer logger.Infof("UpdateProfile end")

	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)
	if err != nil {
		logger.Errorf("failed to get authenticated user UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid request body for profile update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "detail": err.Error()})
		return
	}

	profileInfo, err := h.ProfileService.UpdateProfile(ctx, authID, nickname, &req)
	if err != nil {
		logger.Warnf("failed to update profile for %s: %v", nickname, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
}

func (h *ProfileHandler) GetTheme(c *gin.Context) {
	logger.Infof("GetTheme start")
	defer logger.Infof("GetTheme end")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserID(c)
	if err != nil {
		logger.Errorf("failed to get user UUID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	theme, err := h.ProfileService.GetTheme(ctx, userUUID)
	if err != nil {
		logger.Errorf("failed to get theme for user %s: %v", userUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"theme": theme})
}

// PATCH /profile/theme
func (h *ProfileHandler) SetTheme(c *gin.Context) {
	logger.Infof("SetTheme start")
	defer logger.Infof("SetTheme end")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserID(c)
	if err != nil {
		logger.Errorf("failed to get user UUID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ThemeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "detail": err.Error()})
		return
	}

	theme := req.Theme
	if err := h.ProfileService.SetTheme(ctx, userUUID, theme); err != nil {
		logger.Errorf("failed to set theme for user %s: %v", userUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set theme"})
		return
	}

	c.JSON(http.StatusOK, theme)
}
