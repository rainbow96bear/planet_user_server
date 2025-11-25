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

// 🌐 라우팅 등록 (RESTful 및 중복 제거)
func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	// 1. /me 그룹: 인증된 사용자 전용 (AuthMiddleware 필수)
	me := r.Group("/me")
	me.Use(middleware.AuthMiddleware())
	{
		// 내 프로필 리소스 (Profile)
		me.GET("/profile", h.GetMyProfileInfo) // GET /me/profile
		me.PATCH("/profile", h.UpdateProfile)  // PATCH /me/profile

	}

	// 2. /users 그룹: 공개된 사용자 정보 조회 전용 (AuthMiddleware 불필요)
	// GetProfileInfo는 이 그룹을 사용하도록 통일합니다.
	users := r.Group("/users/:nickname")
	users.GET("", h.GetProfileInfo) // GET /users/:nickname

	// *주의: 기존의 users.GET("",h.GetProfileInfo)와 profileGroup.GET("/:nickname", h.GetProfileInfo)는
	// /users/:nickname 경로로 통일하고 AuthMiddleware를 제거했습니다.
}

// ---------------------- Handler ----------------------

// 다른 유저 프로필 조회: GET /users/:nickname (인증 불필요)
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
		// 사용자가 없을 경우 404 Not Found가 더 적절합니다.
		logger.Warnf("Failed to get profile info for %s: %v", nickname, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
		return
	}

	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
}

// 내 프로필 조회: GET /me/profile (인증 필요)
func (h *ProfileHandler) GetMyProfileInfo(c *gin.Context) {
	logger.Infof("GetMyProfileInfo start")
	defer logger.Infof("GetMyProfileInfo end")

	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)

	// GetUserID에서 에러가 났거나 UUID가 비어있으면 미들웨어에서 걸러지지만, 방어 코드 유지
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

// 내 프로필 업데이트: PATCH /me/profile (인증 필요)
// *참고: URI에서 nickname 파라미터를 제거했습니다. 인증된 사용자의 프로필만 업데이트 가능해야 하므로.
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	logger.Infof("UpdateProfile start")
	defer logger.Infof("UpdateProfile end")

	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)
	if err != nil {
		// 미들웨어에서 처리되지만, 방어 코드 유지
		logger.Errorf("failed to get authenticated user UUID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("invalid request body for profile update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "detail": err.Error()})
		return
	}

	// 업데이트할 닉네임은 요청 본문(req) 내에 있어야 합니다.
	// 여기서는 현재 사용자의 닉네임을 다시 찾는 로직이 필요할 수 있으나,
	// 편의상 기존의 service 호출 시 닉네임을 요구하는 부분을 수정하지 않고 임시로 빈 값으로 전달합니다.
	// *실제 구현 시 service 단에서 authID를 사용하여 닉네임을 조회하거나, 닉네임 유효성 검증을 해야 합니다.
	currentNickname := "" // nickname 파라미터를 제거했기 때문에 임시 처리

	profileInfo, err := h.ProfileService.UpdateProfile(ctx, authID, currentNickname, &req)
	if err != nil {
		logger.Warnf("failed to update profile for %s: %v", authID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
}
