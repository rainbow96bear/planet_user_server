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

type SettingHandler struct {
	SettingService *service.SettingService
}

func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{
		SettingService: settingService,
	}
}

func (h *SettingHandler) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/theme", h.GetTheme)
		userGroup.POST("/theme", h.SetTheme)
	}
}

func (h *SettingHandler) GetTheme(c *gin.Context) {
	logger.Infof("start to get theme")
	defer logger.Infof("end to get theme")

	ctx := c.Request.Context()
	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theme, err := h.SettingService.GetTheme(ctx, authUuid)
	if err != nil {
		logger.Errorf("fail to get theme uuid : %s ERR[%s]", authUuid, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"theme": theme,
	})
}

func (h *SettingHandler) SetTheme(c *gin.Context) {
	logger.Infof("start to set theme")
	defer logger.Infof("end to set theme")

	ctx := c.Request.Context()
	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &dto.Theme{}

	if err := c.ShouldBindJSON(req); err != nil {
		logger.Errorf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Theme != "light" && req.Theme != "dark" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "theme must be 'light' or 'dark'"})
		return
	}

	err = h.SettingService.SetTheme(ctx, authUuid, req.Theme)
	if err != nil {
		logger.Errorf("fail to set theme : %s, uuid : %s ERR[%s]", req.Theme, authUuid, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"theme": req.Theme})
}
