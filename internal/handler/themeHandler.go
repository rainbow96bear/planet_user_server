package handler

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rainbow96bear/planet_user_server/dto"
// 	"github.com/rainbow96bear/planet_user_server/internal/service"
// 	"github.com/rainbow96bear/planet_user_server/middleware"
// 	"github.com/rainbow96bear/planet_user_server/utils"
// 	"github.com/rainbow96bear/planet_utils/pkg/logger"
// )

// type ThemeHandler struct {
// 	ProfileService *service.ProfileService
// }

// func NewThemeHandler(profileService *service.ProfileService) *ThemeHandler {
// 	return &ThemeHandler{
// 		ProfileService: profileService,
// 	}
// }

// // 🌐 라우팅 등록 (RESTful 및 중복 제거)
// func (h *ThemeHandler) RegisterRoutes(r *gin.Engine) {
// 	// 1. /me 그룹: 인증된 사용자 전용 (AccessTokenAuthMiddleware 필수)
// 	me := r.Group("/me")
// 	me.Use(middleware.AccessTokenAuthMiddleware())
// 	{
// 		// 내 테마 설정 리소스 (Theme)
// 		me.GET("/theme", h.GetTheme)   // GET /me/theme
// 		me.PATCH("/theme", h.SetTheme) // PATCH /me/theme
// 	}
// }

// // 내 테마 조회: GET /me/theme (인증 필요)
// func (h *ThemeHandler) GetTheme(c *gin.Context) {
// 	logger.Infof("GetTheme start")
// 	defer logger.Infof("GetTheme end")

// 	ctx := c.Request.Context()
// 	userUUID, err := utils.GetUserID(c)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	theme, err := h.ProfileService.GetTheme(ctx, userUUID)
// 	if err != nil {
// 		logger.Errorf("failed to get theme for user %s: %v", userUUID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get theme"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"theme": theme})
// }

// // 내 테마 설정/업데이트: PATCH /me/theme (인증 필요)
// func (h *ThemeHandler) SetTheme(c *gin.Context) {
// 	logger.Infof("SetTheme start")
// 	defer logger.Infof("SetTheme end")

// 	ctx := c.Request.Context()
// 	userUUID, err := utils.GetUserID(c)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var req dto.ThemeUpdateRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		logger.Warnf("invalid request body: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "detail": err.Error()})
// 		return
// 	}

// 	theme := req.Theme
// 	if err := h.ProfileService.SetTheme(ctx, userUUID, theme); err != nil {
// 		logger.Errorf("failed to set theme for user %s: %v", userUUID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set theme"})
// 		return
// 	}

// 	// 200 OK와 함께 업데이트된 테마 정보를 반환합니다.
// 	c.JSON(http.StatusOK, gin.H{"theme": theme, "message": "Theme updated successfully"})
// }
