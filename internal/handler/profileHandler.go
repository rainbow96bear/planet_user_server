package handler

import (
	"fmt"
	"net/http"
	"planet_utils/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/dto"
	"github.com/rainbow96bear/planet_auth_server/external/oauthClient"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
)

type ProfileHandler struct {
	ProfileService *service.ProfileService
	AuthService *service.AuthService
}

func (h *ProfileHandler)GetProfileInfo(c *gin.Context) {
	logger.Infof("start to get profile")
	defer logger.Infof("end to get profile")
	ctx := c.Request.Context()

	nickname := c.Param("userNickName")
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

	c.JSON(http.StatusOK, gin.H{
		"profile": profileInfo,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	logger.Infof("start to update profile")
	defer logger.Infof("end to update profile")

	ctx := c.Request.Context()

	profileInfo := &dto.ProfileInfo{}

	if err := c.ShouldBindJSON(profileInfo); err != nil {
		logger.Warnf("fail to bind to json about profileInfo : %+v", profileInfo)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.ProfileService.UpdateProfile(ctx, profileInfo); err != nil {
		logger.Warnf("fail to update profileInfo : %+v", profileInfo)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profileInfo,
	})
}