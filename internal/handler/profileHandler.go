package handler

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rainbow96bear/planet_user_server/dto"
// 	"github.com/rainbow96bear/planet_user_server/internal/service"
// 	"github.com/rainbow96bear/planet_user_server/middleware"
// 	"github.com/rainbow96bear/planet_user_server/utils"
// 	"github.com/rainbow96bear/planet_utils/pkg/logger"
// )

// type ProfileHandler struct {
// 	ProfileService *service.ProfileService
// 	FollowService  *service.FollowService
// }

// func NewProfileHandler(profileService *service.ProfileService, followService *service.FollowService) *ProfileHandler {
// 	return &ProfileHandler{
// 		ProfileService: profileService,
// 		FollowService:  followService,
// 	}
// }

// // 🌐 라우팅 등록 (RESTful 및 중복 제거)
// func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
// 	// 1. /me 그룹: 인증된 사용자 전용 (AccessTokenAuthMiddleware 필수)
// 	me := r.Group("/me")
// 	me.Use(middleware.AccessTokenAuthMiddleware())
// 	{
// 		// 내 프로필 리소스 (Profile)
// 		me.GET("/profile", h.GetMyProfileInfo) // GET /me/profile
// 		me.PATCH("/profile", h.UpdateProfile)  // PATCH /me/profile

// 	}

// 	users := r.Group("/users/:nickname")
// 	users.GET("", h.GetProfileInfo)
// }

// // ---------------------- Handler ----------------------

// // 다른 유저 프로필 조회: GET /users/:nickname (인증 불필요)
// func (h *ProfileHandler) GetProfileInfo(c *gin.Context) {
// 	logger.Infof("GetProfileInfo start")
// 	defer logger.Infof("GetProfileInfo end")

// 	ctx := c.Request.Context()
// 	nickname := c.Param("nickname")
// 	if nickname == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
// 		return
// 	}

// 	profileInfo, err := h.ProfileService.GetProfileInfo(ctx, nickname)
// 	if err != nil {
// 		// 사용자가 없을 경우 404 Not Found가 더 적절합니다.
// 		logger.Warnf("Failed to get profile info for %s: %v", nickname, err)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
// }

// // 내 프로필 조회: GET /me/profile (인증 필요)

// // 내 프로필 업데이트: PATCH /me/profile (인증 필요)
// // *참고: URI에서 nickname 파라미터를 제거했습니다. 인증된 사용자의 프로필만 업데이트 가능해야 하므로.
// func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
// 	logger.Infof("UpdateProfile start")
// 	defer logger.Infof("UpdateProfile end")

// 	ctx := c.Request.Context()
// 	authID, err := utils.GetUserID(c)
// 	if err != nil {
// 		// 미들웨어에서 처리되지만, 방어 코드 유지
// 		logger.Errorf("failed to get authenticated user UUID: %v", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var req dto.ProfileUpdateRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		logger.Warnf("invalid request body for profile update: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "detail": err.Error()})
// 		return
// 	}

// 	currentNickname := "" // 닉네임 파라미터가 DTO 내에 포함되었다고 가정하고 임시 처리

// 	profileInfo, err := h.ProfileService.UpdateProfile(ctx, authID, currentNickname, &req)
// 	if err != nil {
// 		logger.Warnf("failed to update profile for %s: %v", authID, err)

// 		// 🌟 핵심 수정: 오류 타입 확인 및 사용자 친화적 응답 🌟
// 		if errors.Is(err, planet_err.ErrNicknameDuplicate) {
// 			// HTTP 409 Conflict 상태 코드 (자원 충돌) 사용
// 			c.JSON(http.StatusConflict, gin.H{
// 				// 사용자에게 보여줄 간결한 메시지
// 				"error": "사용자 이름을 업데이트하지 못했습니다. 해당 사용자 이름은 이미 사용 중일 수 있습니다. 다른 이름을 선택해 주세요.",
// 			})
// 			return
// 		}

// 		// 그 외의 일반적인 오류 처리
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, dto.ToProfileResponse(profileInfo))
// }
